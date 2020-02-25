package kitty

// TODO: Write a function that gets list of objects and.
//=================================
type pages interface {
}

func NewPages(page ...int) pages {
	return "pages instance..."
}

// prototype:
type bookPresenter struct {
	Name string `json:"name"`
}

type bookColl struct {
	pages
	Items []bookPresenter `json:"books"`
}

func NewBookPresenter(book string) bookPresenter {
	return bookPresenter{Name: book}
}

func NewPresenterColl(books []string) bookColl {
	items := make([]bookPresenter, len(books))

	for i, b := range books {
		items[i] = NewBookPresenter(b)
	}

	return bookColl{
		pages: NewPages(1, 2, 3),
		Items: items,
	}
}

//=================================
