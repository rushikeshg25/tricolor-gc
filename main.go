package main

func main() {
	gc := NewGarbageCollector()

	obj1 := gc.Allocate(100)
	obj2 := gc.Allocate(200)
	obj3 := gc.Allocate(150)
	obj4 := gc.Allocate(80)

	//Stack vars, global vars (imp live vars)
	gc.AddRoot(obj1)
	gc.AddRoot(obj2)

	//Adding children refs
	gc.AddReference(obj1, obj3) // obj1 -> obj3
	gc.AddReference(obj2, obj4) // obj2 -> obj4
	gc.AddReference(obj3, obj4) // obj3 -> obj4

}
