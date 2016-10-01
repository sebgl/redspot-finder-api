package main

import (
	"encoding/json"
	"strings"

	"github.com/sebgl/redspot-finder-scraper/scraper"
	"gopkg.in/olivere/elastic.v3"
)

// MaxReturnedItems is the maximum  number of results to be returned
const MaxReturnedItems = 10

// SearchPlaylists returns playlist matching the given query
func SearchPlaylists(query string) ([]scraper.Playlist, error) {
	esQuery := queryTagsFromUserQuery(query).toESQuery()
	results, err := es.Search().
		Index(scraper.Index).
		Type(scraper.DocumentType).
		Query(esQuery).
		Size(MaxReturnedItems).
		Do()
	if err != nil {
		return nil, err
	}
	return PlaylistsFromSearchresults(results)
}

// LastPlaylists returns the newest playlists
func LastPlaylists() ([]scraper.Playlist, error) {
	results, err := es.Search().
		Index(scraper.Index).
		Type(scraper.DocumentType).
		Size(MaxReturnedItems).
		Do()
	if err != nil {
		return nil, err
	}
	return PlaylistsFromSearchresults(results)
}

// PlaylistsFromSearchresults transforms es search results into playlists
func PlaylistsFromSearchresults(results *elastic.SearchResult) ([]scraper.Playlist, error) {
	playlists := make([]scraper.Playlist, len(results.Hits.Hits))
	for i, hit := range results.Hits.Hits {
		var playlist scraper.Playlist
		err := json.Unmarshal(*hit.Source, &playlist)
		if err != nil {
			return nil, err
		}
		playlists[i] = playlist
	}
	return playlists, nil
}

type queryTags struct {
	artists []string
	songs   []string
	query   string
}

func (t queryTags) toESQuery() elastic.Query {
	queries := make([]elastic.Query, 0, 3)
	if len(t.artists) > 0 {
		queries = append(queries,
			elastic.NewTermsQuery("SpotifyData.Tracks.Artists.Name", t.artists))
	}
	if len(t.songs) > 0 {
		queries = append(queries,
			elastic.NewTermsQuery("SpotifyData.Tracks.Name", t.songs))
	}
	if len(t.query) > 0 {
		queries = append(queries,
			elastic.NewQueryStringQuery(t.query))
	}
	return elastic.NewBoolQuery().Must(queries...)
}

func queryTagsFromUserQuery(userQuery string) queryTags {
	tags := queryTags{}
	splits := strings.Split(userQuery, ";")
	for _, s := range splits {
		innerSplits := strings.Split(s, ":")
		if len(innerSplits) < 2 {
			tags.query += s
			continue
		}
		switch innerSplits[0] {
		case "artist:":
			artists := strings.Split(strings.TrimPrefix(s, "artist:"), ",")
			tags.artists = append(tags.artists, artists...)
		case "song:":
			songs := strings.Split(strings.TrimPrefix(s, "song:"), ",")
			tags.songs = append(tags.songs, songs...)
		default:
			tags.query += s
		}
	}
	return tags
}
