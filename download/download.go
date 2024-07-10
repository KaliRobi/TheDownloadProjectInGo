package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
)

type UrlObject struct {
	Id      int
	Content string
}

var waitingGroup = sync.WaitGroup{}

//utility functions

func isValidEndpoint(endpoint string) bool {
	u, err := url.Parse(endpoint)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

func readResponse(response *http.Response) ([]byte, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func reverseSortUrlObjects(UrlObjectList []UrlObject) {
	sort.Slice(UrlObjectList, func(i, j int) bool {
		return UrlObjectList[i].Id > UrlObjectList[j].Id
	})
}

func concatenateContent(objectSlice []UrlObject) string {
	contentSlice := []string{}
	for i := range objectSlice {
		contentSlice = append(contentSlice, objectSlice[i].Content)
	}

	return strings.Join(contentSlice, " ")
}

// for goroutines
func downloadUrl(ctx context.Context, endpoint string, instanceId int, errChannel chan<- error, resultchan chan<- UrlObject) {

	select {
	case <-ctx.Done():
		return
	default:
	}

	if !isValidEndpoint(endpoint) {
		errChannel <- fmt.Errorf("invalid endpoint: %s", endpoint)
		return
	}

	resp, ResponseErr := http.Get(endpoint)

	if ResponseErr != nil {
		select {
		case errChannel <- ResponseErr:
		case <-ctx.Done():
		}
		return
	}
	defer resp.Body.Close()
	content, readErr := readResponse(resp)

	if readErr != nil {
		errChannel <- readErr
		return
	}
	select {
	case resultchan <- UrlObject{Id: instanceId, Content: string(content)}:
	case <-ctx.Done():
		return
	}

	waitingGroup.Done()
}

func ReturnContentOrFail(urlsSlice []string) (string, error) {
	var UrlObjectList = make([]UrlObject, 0, len(urlsSlice))
	ctx, cancel := context.WithCancel(context.Background())
	var errorMessage error

	resultChannel := make(chan UrlObject, len(urlsSlice))
	errorChannel := make(chan error)
	idCounter := 0

	// loop through the slice what contains the URLs start a gorutine for each of them
	for i := range urlsSlice {
		idCounter++
		waitingGroup.Add(1)
		go downloadUrl(ctx, urlsSlice[i], idCounter, errorChannel, resultChannel)

	}

	//monitoring for error thorugh errorChannel and trigger context cancel
	go func() {
		for err := range errorChannel {
			if err != nil {
				errorMessage = fmt.Errorf("error: %v", err)
				cancel()
				return
			}
		}
	}()

	// Processing the result offloading from result channel to the object list
	for i := 0; i < len(urlsSlice); i++ {
		select {
		case res := <-resultChannel:
			UrlObjectList = append(UrlObjectList, res)
		case <-ctx.Done():
			return "", fmt.Errorf("process was cancelled as one of the urls was faulty")
		}
	}

	reverseSortUrlObjects(UrlObjectList)

	appendednConetet := concatenateContent(UrlObjectList)

	if errorMessage != nil {
		return "", errorMessage
	}

	return appendednConetet, nil

}
