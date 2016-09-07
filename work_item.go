package main

//An Identifiable has an ID.
type Identifiable interface {
	ID() int
	SetId(id int)
}

//WorkItem is the representation of a piece of work to be done by an agent, converting a parametrised notebook to static HTML.
type WorkItem struct {
	ItemID       int                    `json:"id"`
	NotebookPath string                 `json:"notebook_path"`
	Parameters   map[string]interface{} `json:"parameters"`
}

//ID returns the ID of the workitem to satisfy the Identifiable interface.
func (wi *WorkItem) ID() int {
	return wi.ItemID
}

//SetId sets the work item's Id
func (wi *WorkItem) SetId(id int) {
	wi.ItemID = id
}
