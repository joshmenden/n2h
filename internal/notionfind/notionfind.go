package notionfind

import (
	"context"
	"fmt"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/jomei/notionapi"
)

func Content(pageTitle, statusProperty, status *string) (*notionapi.Page, []notionapi.Block, error) {
	db, err := findContentDB()
	if err != nil {
		return nil, nil, err
	}

	page, err := findContentPage(db, pageTitle, statusProperty, status)
	if err != nil {
		return nil, nil, err
	}

	blocks, err := getAllPaginatedBlocks(notionapi.BlockID(page.ID))
	if err != nil {
		return nil, nil, err
	}

	return page, blocks, nil
}

func getClient() *notionapi.Client {
	notionAPIKey := config.String("notionAPIKey")
	return notionapi.NewClient(notionapi.Token(notionAPIKey))
}

func findContentDB() (*notionapi.Database, error) {
	dbName := config.String("databaseName")
	searchResults, err := getClient().Search.Do(context.Background(), &notionapi.SearchRequest{Query: dbName})
	if err != nil {
		return nil, err
	}

	if len(searchResults.Results) <= 0 {
		return nil, fmt.Errorf("0 results were found matching: %s", dbName)
	}

	firstResult := searchResults.Results[0]
	if firstResult.GetObject() != notionapi.ObjectTypeDatabase {
		return nil, fmt.Errorf("db query for '%s' doesn't return a database as it's best match", dbName)
	}

	resultDB := firstResult.(*notionapi.Database)
	if !strings.EqualFold(resultDB.Title[0].PlainText, dbName) {
		return nil, fmt.Errorf("hottest match for DB did not match query: %s", dbName)
	}

	return resultDB, nil
}

func findContentPage(db *notionapi.Database, pageTitle, statusProperty, status *string) (*notionapi.Page, error) {
	nameExactFilter := notionapi.PropertyFilter{
		Property: "Name",
		RichText: &notionapi.TextFilterCondition{
			Equals: *pageTitle,
		},
	}

	nameStartsWithFilter := notionapi.PropertyFilter{
		Property: "Name",
		RichText: &notionapi.TextFilterCondition{
			StartsWith: *pageTitle,
		},
	}

	statusFilter := notionapi.PropertyFilter{
		Property: *statusProperty,
		Select: &notionapi.SelectFilterCondition{
			Equals:  *status,
			IsEmpty: false,
		},
	}

	exactFilter := notionapi.AndCompoundFilter{nameExactFilter, statusFilter}
	startsWithFilter := notionapi.AndCompoundFilter{nameStartsWithFilter, statusFilter}

	var dbResponse *notionapi.DatabaseQueryResponse

	// try exact match first
	dbResponse, err := getClient().Database.Query(context.Background(), notionapi.DatabaseID(db.ID), &notionapi.DatabaseQueryRequest{Filter: exactFilter})
	if err != nil {
		return nil, err
	}

	if len(dbResponse.Results) <= 0 {
		// try with starts with filter
		dbResponse, err = getClient().Database.Query(context.Background(), notionapi.DatabaseID(db.ID), &notionapi.DatabaseQueryRequest{Filter: startsWithFilter})
		if err != nil {
			return nil, err
		}
		if len(dbResponse.Results) <= 0 {
			return nil, fmt.Errorf("no pages were found matching the status %s", *status)
		}
	}

	best := dbResponse.Results[0]
	bestTitle := best.Properties["Name"].(*notionapi.TitleProperty)

	if !strings.EqualFold(bestTitle.Title[0].PlainText, *pageTitle) && !strings.HasPrefix(bestTitle.Title[0].PlainText, *pageTitle) {
		return nil, fmt.Errorf("found page does not match given title: %s (found), %s (given)", bestTitle.Title[0].PlainText, *pageTitle)
	}

	return &best, nil
}

func getAllPaginatedBlocks(parentBlockID notionapi.BlockID) (blocks []notionapi.Block, err error) {
	pagination := notionapi.Pagination{
		StartCursor: "",
		PageSize:    100,
	}

	for {
		response, err := getClient().Block.GetChildren(context.Background(), parentBlockID, &pagination)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, response.Results...)

		if !response.HasMore {
			break
		} else {
			pagination.StartCursor = notionapi.Cursor(response.NextCursor)
		}
	}

	return
}
