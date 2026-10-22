package permissions

import "slices"

type Permission string

var (
	ReadUsers   Permission = "users.read"
	CreateUsers Permission = "users.create"
	UpdateUsers Permission = "users.update"
	DeleteUsers Permission = "users.delete"

	ReadDocuments           Permission = "documents.read"
	CreateDocuments         Permission = "documents.create"
	UpdateDocuments         Permission = "documents.update"
	DeleteDocuments         Permission = "documents.delete"
	ReadDocumentsPublicly   Permission = "documents.read-publicly"
	UpdateDocumentsPublicly Permission = "documents.update-publicly"
	DeleteDocumentsPublicly Permission = "documents.delete-publicly"

	CreateLinks Permission = "links.create"
	ReadLinks   Permission = "links.read"
	RevokeLinks Permission = "links.revoke"
)

var servicePermissions = map[string][]Permission{
	"dockloud-api": {
		ReadUsers, CreateUsers, UpdateUsers, DeleteUsers,
		ReadDocuments, CreateDocuments, UpdateDocuments, DeleteDocuments,
	},
	"dockloud-link": {
		CreateLinks, ReadLinks, RevokeLinks,
		ReadDocumentsPublicly, UpdateDocumentsPublicly, DeleteDocumentsPublicly,
	},
}

func (p *Permission) CheckFor(service string) bool {
	permissions, ok := servicePermissions[service]
	if !ok {
		return false
	}
	return slices.Contains(permissions, *p)
}
