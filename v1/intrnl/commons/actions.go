package commons

const (
	ActionNavigation        = "NAVIGATION"
	ActionUpdateParams      = "UPDATE_PARAMS"
	ActionBackNavigation    = "BACK_NAVIGATION"
	ActionForwardNavigation = "FORWARD_NAVIGATION"
	ActionShowBottomSheet   = "SHOW_BOTTOMSHEET"
	ActionCloseBottomSheet  = "CLOSE_BOTTOM_SHEET"
	ActionFetch             = "FETCH"
	ActionAPI               = "API"
	ActionPlayVideo         = "PLAY_VIDEO"
	ActionFetchAPI          = "FETCH_API"
	ActionResumeFlashcards  = "RESUME_FLASHCARDS"
)

const (
	ActionNavigationTypeReplace = "replace"
)

type PopUpImage struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Text  string `json:"alt_text"`
}

type Content struct {
	Heading     string `json:"heading"`
	Description string `json:"description"`
}

type PopupAction struct {
	Type string    `json:"type"`
	Data PopUpData `json:"data"`
}

type PopUpData struct {
	Type            string     `json:"type"`
	Title           string     `json:"title"`
	Image           PopUpImage `json:"image_data"`
	Content         Content    `json:"content"`
	PrimaryCTA      CTA        `json:"primary_cta,omitempty"`
	ShowCloseButton bool       `json:"show_close_button"`
	//SecondaryCTA CTA        `json:"secondary_cta,omitempty"`
}

type PopUp struct {
	Label  string      `json:"label"`
	Action PopupAction `json:"action"`
}

type SecondaryPopUpData struct {
	Type            string     `json:"type"`
	Title           string     `json:"title"`
	Image           PopUpImage `json:"image_data"`
	Content         Content    `json:"content"`
	PrimaryCTA      CTA        `json:"primary_cta,omitempty"`
	SecondaryCTA    CTA        `json:"secondary_cta,omitempty"`
	ShowCloseButton bool       `json:"show_close_button"`
}

type CheckoutPopupAction struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type OpenBottomSheetAction struct {
	Type string              `json:"type"`
	Data OpenBottomSheetData `json:"data"`
}

type OpenBottomSheetData struct {
	Title           string `json:"title"`
	SubText         string `json:"sub_text"`
	ShowCloseButton bool   `json:"show_close_button"`
	PrimaryCTA      CTA    `json:"primary_cta"`
	SecondaryCTA    CTA    `json:"secondary_cta"`
}

type VideoAction struct {
	Type string    `json:"type"`
	Data VideoData `json:"data"`
}

type BannerAction struct {
	Type string           `json:"type"`
	Data BannerActionData `json:"data"`
}

type BannerActionData struct {
	URI    string                 `json:"uri"`
	Type   string                 `json:"type,omitempty"`
	Query  map[string]interface{} `json:"query"`
	Action string                 `json:"action"`
}
