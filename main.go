package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/pkg/errors"
)

type MetaData struct {
	AuthorName    string
	CommitterName string
	SHA           string
	Message       string
}

func (m MetaData) String() string {
	str, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		log.Fatalln("Can't print", err)
		return ""
	}

	return string(str)
}

func (m *MetaData) UnmarshalJSON(buf []byte) error {
	var commit struct {
		SHA    string `json:"sha"`
		Commit struct {
			Author struct {
				Name string `json:"name"`
			} `json:"author"`
			Committer struct {
				Name string `json:"name"`
			} `json:"committer"`
			Message string `json:"message"`
		} `json:"commit"`
	}

	if err := json.Unmarshal(buf, &commit); err != nil {
		return errors.Wrap(err, "parsing into MetaData struct failed")
	}

	m.AuthorName = commit.Commit.Author.Name
	m.CommitterName = commit.Commit.Committer.Name
	m.SHA = commit.SHA
	m.Message = commit.Commit.Message

	return nil
}

type MetaDatas []MetaData

func (ms *MetaDatas) UnmarshalJSON(buf []byte) error {
	// []MetaData is not the same as MetaDatas, and this difference is
	// important!
	var metadatas []MetaData

	if err := json.Unmarshal(buf, &metadatas); err != nil {
		log.Fatalln("error parsing JSON", err)
	}

	// filtering without allocations
	// https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
	cleanedms := metadatas[:0]
	for _, metadata := range metadatas {
		if !strings.HasPrefix(metadata.Message, "WIP") {
			cleanedms = append(cleanedms, metadata)
		}
	}
	*ms = cleanedms

	return nil
}

func main() {
	jsonb, err := ioutil.ReadFile("./commits.json")
	if err != nil {
		log.Fatalln("error reading commits file", err)

	}

	var metadatas1 []MetaData
	if err := json.Unmarshal(jsonb, &metadatas1); err != nil {
		log.Fatalln("error parsing JSON", err)
	}
	fmt.Println(metadatas1)

	var metadatas2 MetaDatas
	if err := json.Unmarshal(jsonb, &metadatas2); err != nil {
		log.Fatalln("error parsing JSON", err)
	}
	fmt.Println(metadatas2)
}
