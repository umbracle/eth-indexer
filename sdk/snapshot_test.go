package sdk

/*
func TestSnapshot(t *testing.T) {
	s1 := &Snapshot1{tree: iradix.New()}
	s1.init()

	obj := s1.get("table-1", "id-1")
	obj.set("a", "b")

	s2 := s1.save()
	s2.init()

	obj2 := s2.get("table-1", "id-2")
	assert.Empty(t, obj2.vals)

	s3 := s2.save()
	s3.init()

	obj1 := s3.get("table-1", "id-1")
	obj1.set("a", "c")
	s3.save()

	fmt.Println(obj1.vals)

	fmt.Println(s1.diffs)
	fmt.Println(s2.diffs)
	fmt.Println(s3.diffs)
}
*/
