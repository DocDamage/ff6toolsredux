// Package docs provides the in-application help system and documentation management.
//
// The docs package manages help topics, categories, and search functionality
// for the application's built-in documentation. It provides structured help
// content organized by categories such as getting started, features,
// troubleshooting, API documentation, and advanced topics.
//
// # Usage
//
//	hs := docs.NewHelpSystem()
//	topic := hs.GetTopic("quickstart")
//	results := hs.Search("inventory")
//
// The help system is designed to be extensible, allowing new topics to be
// registered at runtime for plugin documentation.
package docs
