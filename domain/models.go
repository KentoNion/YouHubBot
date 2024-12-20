package domain

type Link string

type TgChan struct {
	name           string
	link           Link
	sourceChanLink Link
}
