package discourse

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/chadwcarlson/gomeili/config"
	docs "github.com/chadwcarlson/gomeili/index/documents"
	comm "github.com/chadwcarlson/gomeili/parse/discourse/structs"
	"github.com/chadwcarlson/gomeili/utils/ignore"
	req "github.com/chadwcarlson/gomeili/utils/requests"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"strings"
)

// Main method to retrieve Discourse index for Meilisearch.
func Get(p config.Config) docs.Index {

	// Get the categories (i.e. How-to, Tutorials, etc.)
	var categories comm.DiscourseCategories
	body := req.RequestData(p.URL, "categories.json")
	json.Unmarshal(body, &categories)

	// Parse each category individually.
	var allDocuments docs.Index
	io.WriteString(os.Stdout, fmt.Sprintf("\n\033[1mDiscourse API @\033[0m %s\n", p.URL))
	for _, category := range categories.CategoryList.Categories {
		// Check if the category should be ignored, according to `config.yaml`.
		if !ignore.ItemExists(p.Ignore, category.Name) && category.TopicCount > 0 {
			// Get the documents.
			io.WriteString(os.Stdout, fmt.Sprintf("\033[1m %s\033[0m (%v topics)\n", category.Name, category.TopicCount))
			allDocuments = parseCategory(p, category, allDocuments)
		}
	}
	return allDocuments
}

// Parse an indivual category into Meilisearch format.
func parseCategory(p config.Config, category_data comm.Categories, documentsCategory docs.Index) docs.Index {

	// Get the category.
	url_path := "/c/" + category_data.Slug + ".json"
	var category comm.CommunityCategory
	body := req.RequestData(p.URL, url_path)
	json.Unmarshal(body, &category)

	// Get the topic (post) documents for that category.
	documents := parseTopicsPage(p, category, category_data)

	// Add the documents to the index.
	documentsCategory.Documents = append(documentsCategory.Documents, documents.Documents...)

	return documentsCategory

}

// Recursive method for parsing all topics within a category.
// 		The Discourse API is limited to 30-50 posts per API call to a category endpoint.
// 		The API can be adjusted to allow for more, but it will also do so publicly for everyone.
// 		This is the compromise method to get all posts/topics without opening that up completely.
func parseTopicsPage(p config.Config, category comm.CommunityCategory, category_data comm.Categories) docs.Index {

	var documents docs.Index

	// Parse each topic individually.
	bar := progressbar.Default(int64(len(category.TopicList.Topics)))
	for _, topic := range category.TopicList.Topics {

		var document docs.Document

		// Basic fields.
		document.Site = p.Name
		document.Title = topic.Title
		document.Section = category_data.Name
		document.Subsection = category_data.Name

		// URLs.
		rel_url := fmt.Sprintf("/t/%s/%d", topic.Slug, topic.ID)
		full_url := fmt.Sprintf("%s%s", p.URL, rel_url)
		document.URL = full_url
		document.RelativeURL = rel_url

		// DocumentID hash.
		h := sha1.New()
		h.Write([]byte(fmt.Sprintf(full_url)))
		document.DocumentID = fmt.Sprintf("%x", h.Sum(nil))

		// Match `config.yaml` rank, and use React primary/secondary designation.
		document.Rank = p.Rank
		if p.Rank == 1 {
			document.Source = "primary"
		} else {
			document.Source = "secondary"
		}

		// Get the individual post/topic.
		var post comm.CommunityPost
		body := req.RequestData(p.URL, rel_url+".json")
		json.Unmarshal(body, &post)

		// Catch: Testing this on try.discourse.org, sometimes Q&A posts will come back empty.
		// 		I haven't seen the same bug on Community, but this will at least handle that case
		// 		if it does.
		if len(post.PostStream.Posts) > 0 {
			document.Text = post.PostStream.Posts[0].Cooked
			document.Description = post.PostStream.Posts[0].Cooked
		} else {
			document.Text = ""
			document.Description = ""
		}

		// Append the document to the running index.
		documents.Documents = append(documents.Documents, document)

		bar.Add(1)

	}

	// Discourse shows there are more results for a category via `topics_list.more_topics_url`.
	// 		If set, that URL is requested recursively via parseTopicsPage.
	var emptyTopicList comm.TopicList
	if emptyTopicList.MoreTopicsURL != category.TopicList.MoreTopicsURL {
		// Get the next page.
		url_path := "/c/" + category_data.Slug + ".json?" + strings.Split(category.TopicList.MoreTopicsURL, "?")[1]
		var categoryPage comm.CommunityCategory
		body := req.RequestData(p.URL, url_path)
		json.Unmarshal(body, &categoryPage)

		// Get the documents.
		pageDocuments := parseTopicsPage(p, categoryPage, category_data)

		// Append the documents.
		documents.Documents = append(documents.Documents, pageDocuments.Documents...)
	}

	return documents
}
