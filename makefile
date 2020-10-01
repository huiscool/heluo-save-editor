
heluo-save-editor: 
	go build -o heluo-save-editor main.go

read: heluo-save-editor
	./heluo-save-editor

write: heluo-save-editor
	./heluo-save-editor --write

.PHONY : heluo-save-editor