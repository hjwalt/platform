package brave_search

import (
	"net/http"

	"github.com/hjwalt/platform/format"
)

var (
	FailureResponseFormat = format.Json[FailureResponse]()
	SuccessResponseFormat = format.Json[WebSearchResult]()
)

type FailureResponse struct {
	ErrorResponse ErrorResponse `json:"error"`
	Time          string        `json:"time"`
}

func (er FailureResponse) Error() string {
	return er.ErrorResponse.Detail
}

type BraveClient struct {
	Client  *http.Client
	BaseUrl string
}

type Header struct {
	Key   string
	Value string
}

type Param struct {
	Key   string
	Value string
}

// OBJECT COMPONENTS

type ResultContainer[T any] struct {
	Type             string `json:"type"`
	Results          []T    `json:"results"`
	MutatedByGoggles bool   `json:"mutated_by_goggles"`
}

type Profile struct {
	Name     string `json:"name"`
	LongName string `json:"long_name"`
	URL      string `json:"url"`
	Image    string `json:"img"`
}

type MetaURL struct {
	Scheme   string `json:"scheme"`
	NetLoc   string `json:"netloc"`
	Hostname string `json:"hostname"`
	Favicon  string `json:"favicon"`
	Path     string `json:"path"`
}

type Thumbnail struct {
	Src             string `json:"src"`
	Height          int    `json:"height"`
	Width           int    `json:"width"`
	BackgroundColor string `json:"bg_color"`
	Original        string `json:"original"`
	Logo            bool   `json:"logo"`
	Duplicated      bool   `json:"duplicated"`
	Theme           string `json:"theme"`
	Alt             string `json:"alt"`
}

type VideoData struct {
	Duration  string     `json:"duration"`
	Views     int64      `json:"views"`
	Creator   string     `json:"creator"`
	Publisher string     `json:"publisher"`
	Thumbnail *Thumbnail `json:"thumbnail"`
}

type KnowledgeGraphProfile struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         *URL   `json:"url"`
	Thumbnail   *URL   `json:"thumbnail"`
}

type URL struct {
	Original     string        `json:"original"`
	Display      string        `json:"display"`
	Alternatives []string      `json:"alternatives"`
	Canonical    string        `json:"canonical"`
	Mobile       MobileURLItem `json:"mobile"`
}

type MobileURLItem struct {
	Original string `json:"original"`
	AMP      string `json:"amp"`
	Android  string `json:"android"`
	IOS      string `json:"ios"`
}

type Image struct {
	Thumbnail  *Thumbnail       `json:"thumbnail"`
	URL        string           `json:"url"`
	Properties *ImageProperties `json:"properties"`
	Text       string           `json:"text"`
}

type ImageProperties struct {
	URL         string `json:"url"`
	Resized     string `json:"resized"`
	Height      int    `json:"height"`
	Width       int    `json:"width"`
	Format      string `json:"format"`
	ContentSize string `json:"content_size"`
	Placeholder string `json:"placeholder"`
}

type Article struct {
	Author              []Person      `json:"author"`
	Date                string        `json:"date"`
	Publisher           *Organization `json:"publisher"`
	Thumbnail           *Thumbnail    `json:"thumbnail"`
	IsAccessibleForFree bool          `json:"isAccessibleForFree"`
}

type Organization struct {
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Thumbnail *Thumbnail `json:"thumbnail"`
}

type CreativeWork struct {
	Name      string     `json:"name"`
	Thumbnail *Thumbnail `json:"thumbnail"`
	Rating    *Rating    `json:"rating"`
}

type MusicRecording struct {
	Name      string     `json:"name"`
	Thumbnail *Thumbnail `json:"thumbnail"`
	Rating    *Rating    `json:"rating"`
}

type Review struct {
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Thumbnail   *Thumbnail `json:"thumbnail"`
	Description string     `json:"description"`
	Rating      *Rating    `json:"rating"`
}

type Software struct {
	Name           string `json:"name"`
	Author         string `json:"author"`
	Version        string `json:"version"`
	CodeRepository string `json:"codeRepository"`
	Homepage       string `json:"homepage"`
	DatePublished  string `json:"datePublisher"`
	IsNPM          bool   `json:"is_npm"`
	IsPyPi         bool   `json:"is_pypi"`
}

type PostalAddress struct {
	Type            string `json:"type"`
	Country         string `json:"country"`
	PostalCode      string `json:"postalCode"`
	StreetAddress   string `json:"streetAddress"`
	AddressRegion   string `json:"addressRegion"`
	AddressLocality string `json:"addressLocality"`
	DisplayAddress  string `json:"displayAddress"`
}

type OpeningHours struct {
	CurrentDay []DayOpeningHours   `json:"current_day"`
	Days       [][]DayOpeningHours `json:"days"`
}

type DayOpeningHours struct {
	AbbrName string `json:"abbr_name"`
	FullName string `json:"full_name"`
	Opens    string `json:"opens"`
	Closes   string `json:"closes"`
}

type Contact struct {
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}

type Rating struct {
	RatingValue   float32  `json:"ratingValue"`
	BestRating    float32  `json:"bestRating"`
	ReviewCount   int      `json:"reviewCount"`
	Profile       *Profile `json:"profile"`
	IsTripadvisor bool     `json:"is_tripadvisor"`
}

type Recipe struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Thumbnail      *Thumbnail `json:"thumbnail"`
	URL            string     `json:"url"`
	Domain         string     `json:"domain"`
	Favicon        string     `json:"favicon"`
	Time           string     `json:"time"`
	PrepTime       string     `json:"prep_time"`
	CookTime       string     `json:"cook_time"`
	Ingredients    string     `json:"ingredients"`
	Instructions   []HowTo    `json:"instructions"`
	Servings       int        `json:"servings"`
	Calories       int        `json:"calories"`
	Rating         *Rating    `json:"rating"`
	RecipeCategory string     `json:"recipeCategory"`
	RecipeCuisine  string     `json:"recipeCuisine"`
	Video          *VideoData `json:"video"`
}

type HowTo struct {
	Text  string   `json:"text"`
	Name  string   `json:"name"`
	URL   string   `json:"url"`
	Image []string `json:"image"`
}

type DataProvider struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	LongName string `json:"long_name"`
	Image    string `json:"img"`
}

type Unit struct {
	Value float32 `json:"value"`
	Units string  `json:"units"`
}

type Reviews struct {
	Results                  []TripAdvisorReview `json:"results"`
	ViewMoreURL              string              `json:"viewMoreUrl"`
	ReviewsInForeignLanguage bool                `json:"reviews_in_foreign_language"`
}

type TripAdvisorReview struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Rating      *Rating `json:"rating"`
	Author      *Person `json:"author"`
	ReviewURL   string  `json:"review_url"`
	Language    string  `json:"language"`
}

type Person struct {
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	URL       string     `json:"url"`
	Thumbnail *Thumbnail `json:"thumbnail"`
}

type MovieData struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	URL         string     `json:"url"`
	Thumbnail   *Thumbnail `json:"thumbnail"`
	Release     string     `json:"release"`
	Directors   []Person   `json:"directors"`
	Actors      []Person   `json:"actors"`
	Rating      *Rating    `json:"rating"`
	Genres      []string   `json:"genre"`
	Query       string     `json:"query"`
}

type FAQ struct {
	Type    string `json:"type"`
	Results []QA   `json:"results"`
}

type QA struct {
	Question string  `json:"question"`
	Answer   string  `json:"answer"`
	Title    string  `json:"title"`
	URL      string  `json:"url"`
	MetaURL  MetaURL `json:"meta_url"`
}

type QAPage struct {
	Question string  `json:"question"`
	Answer   *Answer `json:"answer"`
}

type Answer struct {
	Text          string `json:"text"`
	Author        string `json:"author"`
	UpvoteCount   int    `json:"upvoteCount"`
	DownvoteCount int    `json:"downvoteCount"`
}

type Book struct {
	Title     string   `json:"title"`
	Author    []Person `json:"author"`
	Date      string   `json:"date"`
	Price     *Price   `json:"price"`
	Pages     int      `json:"pages"`
	Publisher *Person  `json:"publisher"`
	Rating    *Rating  `json:"rating"`
}

type Price struct {
	Price         string `json:"price"`
	PriceCurrency string `json:"price_currency"`
}

type Product struct {
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Price       string     `json:"price"`
	Thumbnail   *Thumbnail `json:"thumbnail"`
	Description string     `json:"description"`
	Offers      []Offer    `json:"offers"`
	Rating      *Rating    `json:"rating"`
}

type Offer struct {
	URL           string `json:"url"`
	Price         string `json:"price"`
	PriceCurrency string `json:"priceCurrency"`
}

type ForumData struct {
	ForumName  string `json:"forum_name"`
	NumAnswers int    `json:"num_answers"`
	Score      string `json:"score"`
	Question   string `json:"question"`
	TopComment string `json:"top_comment"`
}

type Mixed struct {
	Type string `json:"type"`
	Main []ResultReference
	Top  []ResultReference
	Side []ResultReference
}

type ResultReference struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	All   bool   `json:"all"`
}

type Query struct {
	Original             string    `json:"original"`
	ShowStrictWarning    bool      `json:"show_strict_warning"`
	Altered              string    `json:"altered"`
	Safesearch           bool      `json:"safesearch"`
	IsNavigational       bool      `json:"is_navigational"`
	IsGeolocal           bool      `json:"is_geolocal"`
	LocalDecision        string    `json:"local_decision"`
	LocalLocationsIdx    int       `json:"local_locations_idx"`
	IsTrending           bool      `json:"is_trending"`
	IsNewsBreaking       bool      `json:"is_news_breaking"`
	AskForLocation       bool      `json:"ask_for_location"`
	Language             *Language `json:"language"`
	SpellcheckOff        bool      `json:"spellcheck_off"`
	Country              string    `json:"country"`
	BadResults           bool      `json:"bad_results"`
	ShouldFallback       bool      `json:"should_fallback"`
	Lat                  string    `json:"lat"`
	Long                 string    `json:"long"`
	PostalCode           string    `json:"postal_code"`
	City                 string    `json:"city"`
	State                string    `json:"state"`
	HeaderCountry        string    `json:"header_country"`
	MoreResultsAvailable bool      `json:"more_results_available"`
	CustomLocationLabel  string    `json:"custom_location_label"`
	RedditCluster        string    `json:"reddit_cluster"`
	SummaryKey           string    `json:"summary_key"`
}

type Language struct {
	Main string `json:"main"`
}

type Summarizer struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type ErrorMeta struct {
	Component string           `json:"component"`
	Errors    []ErrorMetaError `json:"errors"`
}

type ErrorMetaError struct {
	Loc     []string     `json:"loc"`
	Message string       `json:"msg"`
	Type    string       `json:"type"`
	Context ErrorContext `json:"ctx"`
	Input   string       `json:"input"`
}

type ErrorContext struct {
	EnumValues []string `json:"enum_values"`
}

type SpellcheckResultItem struct {
	Query string `json:"query"`
}

// RESULTS

type Result struct {
	Title          string   `json:"title"`
	URL            string   `json:"url"`
	IsSourceLocal  bool     `json:"is_source_local"`
	IsSourceBoth   bool     `json:"is_source_both"`
	Description    string   `json:"description"`
	PageAge        string   `json:"page_age"`
	PageFetched    string   `json:"page_fetched"`
	Profile        *Profile `json:"profile"`
	Language       string   `json:"language"`
	FamilyFriendly bool     `json:"family_friendly"`
}

type NewsResult struct {
	Result
	MetaURL   MetaURL    `json:"meta_url"`
	Source    string     `json:"source"`
	Breaking  bool       `json:"breaking"`
	Thumbnail *Thumbnail `json:"thumbnail"`
	Age       string     `json:"age"`
}

type VideoResult struct {
	Result
	Type      string     `json:"type"`
	Data      *VideoData `json:"video"`
	MetaURL   MetaURL    `json:"meta_url"`
	Thumbnail *Thumbnail `json:"thumbnail"`
	Age       string     `json:"age"`
}

type SearchResult struct {
	Result
	Type          string      `json:"type"`
	DeepResults   *DeepResult `json:"deep_results"`
	Schemas       any         `json:"schemas"`
	MetaURL       MetaURL     `json:"meta_url"`
	Thumbnail     *Thumbnail  `json:"thumbnail"`
	Age           string      `json:"age"`
	Language      string      `json:"language"`
	ContentType   string      `json:"content_type"`
	ExtraSnippets []string    `json:"extra_snippets"`

	Subtype        string          `json:"subtype"`
	Article        *Article        `json:"article"`
	Book           *Book           `json:"book"`
	Cluster        []Result        `json:"cluster"`
	ClusterType    string          `json:"cluster_type"`
	CreativeWork   *CreativeWork   `json:"creative_work"`
	FAQ            *FAQ            `json:"faq"`
	Location       *LocationResult `json:"location"`
	Movie          *MovieData      `json:"movie"`
	MusicRecording *MusicRecording `json:"music_recording"`
	ProductCluster []Product       `json:"product_cluster"`
	QA             *QAPage         `json:"qa"`
	Rating         *Rating         `json:"rating"`
	Recipe         *Recipe         `json:"recipe"`
	Review         *Review         `json:"review"`
	Software       *Software       `json:"software"`
	Video          *VideoData      `json:"video"`
}

type DeepResult struct {
	News    []NewsResult            `json:"news"`
	Buttons []ButtonResult          `json:"buttons"`
	Social  []KnowledgeGraphProfile `json:"social"`
	Videos  []VideoResult           `json:"videos"`
	Images  []Image                 `json:"images"`
}

type ButtonResult struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type LocationResult struct {
	Result

	Type           string          `json:"type"`
	ProviderURL    string          `json:"provider_url"`
	Coordinates    []float32       `json:"coordinates"`
	ZoomLevel      int             `json:"zoom_level"`
	Thumbnail      *Thumbnail      `json:"thumbnail"`
	PostalAddress  *PostalAddress  `json:"postal_address"`
	OpeningHours   *OpeningHours   `json:"opening_hours"`
	Contact        *Contact        `json:"contact"`
	PriceRange     string          `json:"price_range"`
	Rating         *Rating         `json:"rating"`
	Distance       *Unit           `json:"distance"`
	Profiles       []DataProvider  `json:"profiles"`
	Reviews        *Reviews        `json:"reviews"`
	Pictures       *PictureResults `json:"pictures"`
	ServesCuisine  []string        `json:"serves_cuisine"`
	Timezone       string          `json:"timezone"`
	TimezoneOffset float32         `json:"timezone_offset"`
	Categories     []string        `json:"categories"`
	IconCategory   string          `json:"icon_category"`
}

type DiscussionResult struct {
	SearchResult

	Type string    `json:"type"`
	Data ForumData `json:"data"`
}

type GraphInfoBox struct {
	Result

	Type            string          `json:"type"`
	Position        int             `json:"position"`
	Label           string          `json:"label"`
	Category        string          `json:"category"`
	LongDesc        string          `json:"long_desc"`
	Thumbnail       *Thumbnail      `json:"thumbnail"`
	Attributes      [][]string      `json:"attributes"`
	Profiles        []Profile       `json:"profiles"`
	WebsiteURL      string          `json:"website_url"`
	AttributesShown int             `json:"attributes_shown"`
	Ratings         []Rating        `json:"ratings"`
	Providers       []DataProvider  `json:"providers"`
	Distance        *Unit           `json:"distance"`
	Images          []Thumbnail     `json:"images"`
	Movie           *MovieData      `json:"movie"`
	Data            *QAPage         `json:"data"`
	FoundInURLs     []string        `json:"found_in_urls"`
	MetaURL         MetaURL         `json:"meta_url"`
	Location        *LocationResult `json:"location"`
	Coordinates     []float32       `json:"coordinates"`
}

type PictureResults struct {
	Results     []Thumbnail `json:"results"`
	ViewMoreURL string      `json:"viewMoreUrl"`
}

type ImageResult struct {
	Type        string           `json:"type"`
	Title       string           `json:"title"`
	URL         string           `json:"url"`
	Source      string           `json:"source"`
	PageFetched string           `json:"page_fetched"`
	Thumbnail   *Thumbnail       `json:"thumbnail"`
	Properties  *ImageProperties `json:"properties"`
	MetaURL     *MetaURL         `json:"meta_url"`
}

// AGGREGATE

type ErrorResponse struct {
	ID       string    `json:"id"`
	Status   int       `json:"status"`
	Code     string    `json:"code"`
	Detail   string    `json:"detail"`
	Meta     ErrorMeta `json:"meta"`
	RawQuery string    `json:"raw_query"`
	Time     string    `json:"-"`
}

func (er ErrorResponse) Error() string {
	return er.Detail
}

type WebSearchResult struct {
	Type        string                             `json:"type"`
	Discussions *ResultContainer[DiscussionResult] `json:"discussions"`
	FAQ         *ResultContainer[QA]               `json:"faq"`
	InfoBox     *ResultContainer[GraphInfoBox]     `json:"infobox"`
	Locations   *ResultContainer[LocationResult]   `json:"locations"`
	Mixed       *Mixed                             `json:"mixed"`
	News        *ResultContainer[NewsResult]       `json:"news"`
	Query       *Query                             `json:"query"`
	Videos      *ResultContainer[VideoResult]      `json:"videos"`
	Web         *ResultContainer[SearchResult]     `json:"web"`
	Summarizer  *Summarizer                        `json:"summarizer"`
}

type ImageSearchResult struct {
	ResultContainer[ImageResult]
	Query *Query `json:"query"`
}

type VideoSearchResult struct {
	ResultContainer[VideoResult]
	Query *Query `json:"query"`
}

type SpellcheckResult struct {
	Type    string                 `json:"type"`
	Query   *Query                 `json:"query"`
	Results []SpellcheckResultItem `json:"results"`
}
