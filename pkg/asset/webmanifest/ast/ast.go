package ast

import (
	"net/url"
)

type (
	WebManifest struct {
		Dir                       string
		Lang                      string
		ShortName                 string
		Description               string
		Icons                     []*ImageResource
		Screenshots               []*ImageResource
		Categories                []string
		IARCRatingID              string
		StartURL                  *url.URL
		Display                   string
		Orientation               string
		ThemeColor                string
		BackgroundColor           string
		Scope                     string
		ServiceWorker             []*ServiceWorkerRegistration
		RelatedApplications       []*ExternalApplicationResource
		PreferRelatedApplications bool
	}

	ImageResource struct {
		Source   *url.URL
		Sizes    []string
		Type     string
		Purpose  string
		Platform string
	}

	ExternalApplicationResource struct {
		Platform       string
		URL            *url.URL
		ID             string
		MinimumVersion string
		Fingerprints   []*Fingerprint
	}

	ServiceWorkerRegistration struct {
	}

	Fingerprint struct {
		Type, Value string
	}
)
