package utils

// need to convert to generic type
type Queue struct {
	Elements []string
}

func (q *Queue) Enqueue(elem string) {
	q.Elements = append(q.Elements, elem)
}

func (q *Queue) Dequeue() string {
	if q.IsEmpty() {
		return ""
	}

	element := q.Elements[0]

	if q.GetLength() == 1 {
		q.Elements = nil
		return element
	}

	q.Elements = q.Elements[1:]

	return element
}

func (q *Queue) GetLength() int {
	return len(q.Elements)
}

func (q *Queue) IsEmpty() bool {
	return len(q.Elements) == 0
}
