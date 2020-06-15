
package structs

import (
	"time"
)


// All Categories
type DiscourseCategories struct {
	CategoryList CategoryList `json:"category_list"`
}
type Categories struct {
	ID                           int         `json:"id"`
	Name                         string      `json:"name"`
	Color                        string      `json:"color"`
	TextColor                    string      `json:"text_color"`
	Slug                         string      `json:"slug"`
	TopicCount                   int         `json:"topic_count"`
	PostCount                    int         `json:"post_count"`
	Position                     int         `json:"position"`
	Description                  string      `json:"description"`
	DescriptionText              string      `json:"description_text"`
	TopicURL                     string      `json:"topic_url"`
	ReadRestricted               bool        `json:"read_restricted"`
	Permission                   interface{} `json:"permission"`
	NotificationLevel            interface{} `json:"notification_level"`
	TopicTemplate                string      `json:"topic_template"`
	HasChildren                  bool        `json:"has_children"`
	SortOrder                    string      `json:"sort_order"`
	SortAscending                interface{} `json:"sort_ascending"`
	ShowSubcategoryList          bool        `json:"show_subcategory_list"`
	NumFeaturedTopics            int         `json:"num_featured_topics"`
	DefaultView                  string      `json:"default_view"`
	SubcategoryListStyle         string      `json:"subcategory_list_style"`
	DefaultTopPeriod             string      `json:"default_top_period"`
	MinimumRequiredTags          int         `json:"minimum_required_tags"`
	NavigateToFirstPostAfterRead bool        `json:"navigate_to_first_post_after_read"`
	TopicsDay                    int         `json:"topics_day"`
	TopicsWeek                   int         `json:"topics_week"`
	TopicsMonth                  int         `json:"topics_month"`
	TopicsYear                   int         `json:"topics_year"`
	TopicsAllTime                int         `json:"topics_all_time"`
	DescriptionExcerpt           string      `json:"description_excerpt"`
	UploadedLogo                 interface{} `json:"uploaded_logo"`
	UploadedBackground           interface{} `json:"uploaded_background"`
}
type CategoryList struct {
	CanCreateCategory bool         `json:"can_create_category"`
	CanCreateTopic    bool         `json:"can_create_topic"`
	Draft             interface{}  `json:"draft"`
	DraftKey          string       `json:"draft_key"`
	DraftSequence     interface{}  `json:"draft_sequence"`
	Categories        []Categories `json:"categories"`
}

// Single Category
type CommunityCategory struct {
	Users         []Users         `json:"users"`
	PrimaryGroups []PrimaryGroups `json:"primary_groups"`
	TopicList     TopicList       `json:"topic_list"`
}
type Users struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	AvatarTemplate string `json:"avatar_template"`
}
type PrimaryGroups struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	FlairURL     string `json:"flair_url"`
	FlairBgColor string `json:"flair_bg_color"`
	FlairColor   string `json:"flair_color"`
}
type Posters struct {
	Extras         string `json:"extras"`
	Description    string `json:"description"`
	UserID         int    `json:"user_id"`
	PrimaryGroupID int    `json:"primary_group_id"`
}
type Topics struct {
	ID                 int           `json:"id"`
	Title              string        `json:"title"`
	FancyTitle         string        `json:"fancy_title"`
	Slug               string        `json:"slug"`
	PostsCount         int           `json:"posts_count"`
	ReplyCount         int           `json:"reply_count"`
	HighestPostNumber  int           `json:"highest_post_number"`
	ImageURL           string        `json:"image_url"`
	CreatedAt          time.Time     `json:"created_at"`
	LastPostedAt       time.Time     `json:"last_posted_at"`
	Bumped             bool          `json:"bumped"`
	BumpedAt           time.Time     `json:"bumped_at"`
	Unseen             bool          `json:"unseen"`
	Pinned             bool          `json:"pinned"`
	Unpinned           interface{}   `json:"unpinned"`
	Visible            bool          `json:"visible"`
	Closed             bool          `json:"closed"`
	Archived           bool          `json:"archived"`
	Bookmarked         interface{}   `json:"bookmarked"`
	Liked              interface{}   `json:"liked"`
	Tags               []interface{} `json:"tags"`
	Views              int           `json:"views"`
	LikeCount          int           `json:"like_count"`
	HasSummary         bool          `json:"has_summary"`
	Archetype          string        `json:"archetype"`
	LastPosterUsername string        `json:"last_poster_username"`
	CategoryID         int           `json:"category_id"`
	PinnedGlobally     bool          `json:"pinned_globally"`
	FeaturedLink       interface{}   `json:"featured_link"`
	Posters            []Posters     `json:"posters"`
}
type TopicList struct {
	MoreTopicsURL  string 		 `json:"more_topics_url"`
	CanCreateTopic bool        `json:"can_create_topic"`
	Draft          interface{} `json:"draft"`
	DraftKey       string      `json:"draft_key"`
	DraftSequence  interface{} `json:"draft_sequence"`
	PerPage        int         `json:"per_page"`
	Topics         []Topics    `json:"topics"`
}

// Single Post

type CommunityPost struct {
	PostStream struct {
		Posts []struct {
			ID                       int         `json:"id"`
			Name                     string      `json:"name"`
			Username                 string      `json:"username"`
			AvatarTemplate           string      `json:"avatar_template"`
			CreatedAt                time.Time   `json:"created_at"`
			Cooked                   string      `json:"cooked"`
			PostNumber               int         `json:"post_number"`
			PostType                 int         `json:"post_type"`
			UpdatedAt                time.Time   `json:"updated_at"`
			ReplyCount               int         `json:"reply_count"`
			ReplyToPostNumber        interface{} `json:"reply_to_post_number"`
			QuoteCount               int         `json:"quote_count"`
			IncomingLinkCount        int         `json:"incoming_link_count"`
			Reads                    int         `json:"reads"`
			Score                    float64     `json:"score"`
			Yours                    bool        `json:"yours"`
			TopicID                  int         `json:"topic_id"`
			TopicSlug                string      `json:"topic_slug"`
			DisplayUsername          string      `json:"display_username"`
			PrimaryGroupName         string      `json:"primary_group_name"`
			PrimaryGroupFlairURL     string      `json:"primary_group_flair_url"`
			PrimaryGroupFlairBgColor string      `json:"primary_group_flair_bg_color"`
			PrimaryGroupFlairColor   string      `json:"primary_group_flair_color"`
			Version                  int         `json:"version"`
			CanEdit                  bool        `json:"can_edit"`
			CanDelete                bool        `json:"can_delete"`
			CanRecover               bool        `json:"can_recover"`
			CanWiki                  bool        `json:"can_wiki"`
			LinkCounts               []struct {
				URL        string `json:"url"`
				Internal   bool   `json:"internal"`
				Reflection bool   `json:"reflection"`
				Title      string `json:"title,omitempty"`
				Clicks     int    `json:"clicks"`
			} `json:"link_counts"`
			Read               bool          `json:"read"`
			UserTitle          string        `json:"user_title"`
			ActionsSummary     []interface{} `json:"actions_summary"`
			Moderator          bool          `json:"moderator"`
			Admin              bool          `json:"admin"`
			Staff              bool          `json:"staff"`
			UserID             int           `json:"user_id"`
			Hidden             bool          `json:"hidden"`
			TrustLevel         int           `json:"trust_level"`
			DeletedAt          interface{}   `json:"deleted_at"`
			UserDeleted        bool          `json:"user_deleted"`
			EditReason         interface{}   `json:"edit_reason"`
			CanViewEditHistory bool          `json:"can_view_edit_history"`
			Wiki               bool          `json:"wiki"`
		} `json:"posts"`
		Stream []int `json:"stream"`
	} `json:"post_stream"`
	TimelineLookup  [][]int `json:"timeline_lookup"`
	SuggestedTopics []struct {
		ID                int         `json:"id"`
		Title             string      `json:"title"`
		FancyTitle        string      `json:"fancy_title"`
		Slug              string      `json:"slug"`
		PostsCount        int         `json:"posts_count"`
		ReplyCount        int         `json:"reply_count"`
		HighestPostNumber int         `json:"highest_post_number"`
		ImageURL          string      `json:"image_url"`
		CreatedAt         time.Time   `json:"created_at"`
		LastPostedAt      time.Time   `json:"last_posted_at"`
		Bumped            bool        `json:"bumped"`
		BumpedAt          time.Time   `json:"bumped_at"`
		Unseen            bool        `json:"unseen"`
		Pinned            bool        `json:"pinned"`
		Unpinned          interface{} `json:"unpinned"`
		Visible           bool        `json:"visible"`
		Closed            bool        `json:"closed"`
		Archived          bool        `json:"archived"`
		Bookmarked        interface{} `json:"bookmarked"`
		Liked             interface{} `json:"liked"`
		Tags              []string    `json:"tags"`
		Archetype         string      `json:"archetype"`
		LikeCount         int         `json:"like_count"`
		Views             int         `json:"views"`
		CategoryID        int         `json:"category_id"`
		FeaturedLink      interface{} `json:"featured_link"`
		Posters           []struct {
			Extras      string `json:"extras"`
			Description string `json:"description"`
			User        struct {
				ID             int    `json:"id"`
				Username       string `json:"username"`
				Name           string `json:"name"`
				AvatarTemplate string `json:"avatar_template"`
			} `json:"user"`
		} `json:"posters"`
	} `json:"suggested_topics"`
	Tags              []string    `json:"tags"`
	ID                int         `json:"id"`
	Title             string      `json:"title"`
	FancyTitle        string      `json:"fancy_title"`
	PostsCount        int         `json:"posts_count"`
	CreatedAt         time.Time   `json:"created_at"`
	Views             int         `json:"views"`
	ReplyCount        int         `json:"reply_count"`
	LikeCount         int         `json:"like_count"`
	LastPostedAt      time.Time   `json:"last_posted_at"`
	Visible           bool        `json:"visible"`
	Closed            bool        `json:"closed"`
	Archived          bool        `json:"archived"`
	HasSummary        bool        `json:"has_summary"`
	Archetype         string      `json:"archetype"`
	Slug              string      `json:"slug"`
	CategoryID        int         `json:"category_id"`
	WordCount         int         `json:"word_count"`
	DeletedAt         interface{} `json:"deleted_at"`
	UserID            int         `json:"user_id"`
	FeaturedLink      interface{} `json:"featured_link"`
	PinnedGlobally    bool        `json:"pinned_globally"`
	PinnedAt          interface{} `json:"pinned_at"`
	PinnedUntil       interface{} `json:"pinned_until"`
	Draft             interface{} `json:"draft"`
	DraftKey          string      `json:"draft_key"`
	DraftSequence     interface{} `json:"draft_sequence"`
	Unpinned          interface{} `json:"unpinned"`
	Pinned            bool        `json:"pinned"`
	CurrentPostNumber int         `json:"current_post_number"`
	HighestPostNumber int         `json:"highest_post_number"`
	DeletedBy         interface{} `json:"deleted_by"`
	ActionsSummary    []struct {
		ID     int  `json:"id"`
		Count  int  `json:"count"`
		Hidden bool `json:"hidden"`
		CanAct bool `json:"can_act"`
	} `json:"actions_summary"`
	ChunkSize        int         `json:"chunk_size"`
	Bookmarked       interface{} `json:"bookmarked"`
	TopicTimer       interface{} `json:"topic_timer"`
	MessageBusLastID int         `json:"message_bus_last_id"`
	ParticipantCount int         `json:"participant_count"`
	Details          struct {
		NotificationLevel int `json:"notification_level"`
		Participants      []struct {
			ID                       int    `json:"id"`
			Username                 string `json:"username"`
			Name                     string `json:"name"`
			AvatarTemplate           string `json:"avatar_template"`
			PostCount                int    `json:"post_count"`
			PrimaryGroupName         string `json:"primary_group_name"`
			PrimaryGroupFlairURL     string `json:"primary_group_flair_url"`
			PrimaryGroupFlairColor   string `json:"primary_group_flair_color"`
			PrimaryGroupFlairBgColor string `json:"primary_group_flair_bg_color"`
		} `json:"participants"`
		CreatedBy struct {
			ID             int    `json:"id"`
			Username       string `json:"username"`
			Name           string `json:"name"`
			AvatarTemplate string `json:"avatar_template"`
		} `json:"created_by"`
		LastPoster struct {
			ID             int    `json:"id"`
			Username       string `json:"username"`
			Name           string `json:"name"`
			AvatarTemplate string `json:"avatar_template"`
		} `json:"last_poster"`
		Links []struct {
			URL        string `json:"url"`
			Title      string `json:"title"`
			Internal   bool   `json:"internal"`
			Attachment bool   `json:"attachment"`
			Reflection bool   `json:"reflection"`
			Clicks     int    `json:"clicks"`
			UserID     int    `json:"user_id"`
			Domain     string `json:"domain"`
			RootDomain string `json:"root_domain"`
		} `json:"links"`
	} `json:"details"`
}
