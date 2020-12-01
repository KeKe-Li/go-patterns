package main

type Employee struct {
	Name     string
	Age      int
	Vacation int
	Salary   int
}

var list = []Employee{
	{"keke", 20, 0, 100000},
	{"jame", 30, 0, 100000},
	{"Array", 36, 0, 400000},
	{"Jack", 33, 10, 300000},
	{"Mike", 32, 20, 600000},
}

func EmployeeCountIf(list []Employee, fn func(e *Employee) bool) int {
	count := 0
	for i, _ := range list {
		if fn(&list[i]) {
			count += 1
		}
	}
	return count
}

func EmployeeFilterIn(list []Employee, fn func(e *Employee) bool) []Employee {
	var newList []Employee
	for i, _ := range list {
		if fn(&list[i]) {
			newList = append(newList, list[i])
		}
	}
	return newList
}

func EmployeeSumIf(list []Employee, fn func(e *Employee) int) int {
	var sum = 0
	for i, _ := range list {
		sum += fn(&list[i])
	}
	return sum
}

//younger_pay := EmployeeSumIf(list, func(e *Employee) int { if e.Age < 30 {
//	return e.Salary } else {
//	return 0 }
//})
//total_pay := EmployeeSumIf(list, func(e *Employee) int {
//	return e.Salary
//})
//fmt.Printf("Total Salary: %d\n", total_pay) //Total Salary: 43500
