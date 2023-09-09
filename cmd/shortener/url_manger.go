package main

import "strings"

type urlManager struct {
	scheme  string
	host    string
	baseURI string
}

func newUrlManager(scheme string, host string, baseURI string) *urlManager {
	return &urlManager{
		scheme:  scheme,
		host:    host,
		baseURI: baseURI,
	}
}

func (receiver urlManager) buildFullURLByShortcut(shortcut string) string {
	return receiver.getBaseURL() + shortcut
}

func (receiver urlManager) getShortcutFromFullURL(url string) string {
	return strings.TrimPrefix(url, receiver.getBaseURL())
}

func (receiver urlManager) getShortcutFromURI(url string) string {
	return strings.TrimPrefix(url, receiver.baseURI)
}

func (receiver urlManager) getBaseURL() string {
	return receiver.scheme + "://" + receiver.host + receiver.baseURI
}
