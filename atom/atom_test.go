package atom

import (
	"testing"
	"time"
)

func areEqual[T comparable](t *testing.T, expected, actual T) {
	if expected != actual {
		t.Errorf("Expected: %v, Actual: %v", expected, actual)
	}
}

func Test_AtomFeed(t *testing.T) {
	pubTime, _ := time.Parse(time.RFC3339, "2022-10-27T22:07:58+01:00")
	feed := NewFeed(
		"https://example.org/feed.atom",
		NewText("text", "Example Feed"),
		pubTime).
		SetAuthor(
			NewPerson("John Doe").
				SetEmail("john.doe@example.org").
				SetURI("https://example.org/john-doe")).
		SetGenerator(
			NewGenerator("Unit Test").
				SetURI("https://example.org/unit-test").
				SetVersion("v1")).
		SetRights(
			NewText("text", "Copyright © 2021 Example.org")).
		AddEntry(
			NewEntry(
				"https://example.org/entry/1",
				NewText("text", "Entry 1"),
				pubTime).
				SetContent(NewText("text", "Entry 1 Content"))).
		AddEntry(
			NewEntry(
				"https://example.org/entry/2",
				NewText("text", "Entry 2"),
				pubTime).
				SetContent(NewText("text", "Entry 2 Content")))

	bytes, err := feed.ToXML(true, true)
	if err != nil {
		t.Fatal(err)
	}

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <id>https://example.org/feed.atom</id>
  <title type="text">Example Feed</title>
  <updated>2022-10-27T22:07:58+01:00</updated>
  <author>
    <name>John Doe</name>
    <email>john.doe@example.org</email>
    <uri>https://example.org/john-doe</uri>
  </author>
  <generator uri="https://example.org/unit-test" version="v1">Unit Test</generator>
  <rights type="text">Copyright © 2021 Example.org</rights>
  <entry>
    <id>https://example.org/entry/1</id>
    <title type="text">Entry 1</title>
    <updated>2022-10-27T22:07:58+01:00</updated>
    <content type="text">Entry 1 Content</content>
  </entry>
  <entry>
    <id>https://example.org/entry/2</id>
    <title type="text">Entry 2</title>
    <updated>2022-10-27T22:07:58+01:00</updated>
    <content type="text">Entry 2 Content</content>
  </entry>
</feed>`

	areEqual(t, expected, string(bytes))
}
