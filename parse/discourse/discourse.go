
package discourse

import (
    "fmt"
    "io"
    "os"
    "encoding/json"
    "crypto/sha1"
    "github.com/schollz/progressbar/v3"
    "github.com/chadwcarlson/gomeili/config"
    "github.com/chadwcarlson/gomeili/utils/ignore"
    req "github.com/chadwcarlson/gomeili/utils/requests"
    docs "github.com/chadwcarlson/gomeili/index/documents"
    comm "github.com/chadwcarlson/gomeili/parse/discourse/structs"
)

func Get(p config.Config) docs.Index {

  // Get the categories.
  var categories comm.DiscourseCategories
  body := req.RequestData(p.URL, "categories.json")
  json.Unmarshal(body, &categories)

  // Parse each category individually.
  var allDocuments docs.Index
  io.WriteString(os.Stdout, fmt.Sprintf("\n\033[1mDiscourse API @\033[0m %s\n", p.URL))
  for _, category := range categories.CategoryList.Categories {
    if !ignore.ItemExists(p.Ignore, category.Name) && category.TopicCount > 0 {
      io.WriteString(os.Stdout, fmt.Sprintf("\033[1m %s\033[0m (%v topics)\n", category.Name, category.TopicCount))
      allDocuments = parseCategory(p, category, allDocuments)
    }
  }

  return allDocuments

}


func parseCategory(p config.Config, category_data comm.Categories, documentsCategory docs.Index) docs.Index {

  // Get the category.
  url_path := "/c/" + category_data.Slug + ".json"
  var category comm.CommunityCategory
  body := req.RequestData(p.URL, url_path)
  json.Unmarshal(body, &category)

  // Parse each topic individually.
  bar := progressbar.Default(int64(len(category.TopicList.Topics)))
  for _, topic := range category.TopicList.Topics {

    var document docs.Document

    document.Site = p.Name
    document.Section = category_data.Name
    document.Title = topic.Title

    rel_url := fmt.Sprintf("/t/%s/%d", topic.Slug, topic.ID)
    full_url := fmt.Sprintf("%s%s", p.URL, rel_url)
    document.URL = full_url
    document.RelativeURL = rel_url

    h := sha1.New()
    h.Write([]byte(fmt.Sprintf(full_url)))
    document.DocumentID = fmt.Sprintf("%x", h.Sum(nil))

    document.Rank = p.Rank
    if p.Rank == 1 {
      document.Source = "primary"
    } else {
      document.Source = "secondary"
    }

    var post comm.CommunityPost
    body := req.RequestData(p.URL, rel_url + ".json")
    json.Unmarshal(body, &post)

    // Catch: Testing this on try.discourse.org, sometimes Q&A posts will come back empty.
    if len(post.PostStream.Posts) > 0 {
      document.Text = post.PostStream.Posts[0].Cooked
    } else {
      document.Text = ""
    }
    document.Subsection = ""

    documentsCategory.Documents = append(documentsCategory.Documents, document)

    bar.Add(1)

  }

  return documentsCategory


}
