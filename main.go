package main

import (
	"database/sql"
	"fmt"

	"github.com/k0kubun/pp"
	_ "github.com/lib/pq"
)

type Student struct {
	Id   int
	Name string
	Age  int
}

type Course struct {
	Id         int
	CourseName string
	Price      int
}

type StudentCourse struct {
	Id        int
	StudentId Student
	CourseId  Course
}

type ReceivedData struct {
	Id         int
	Name       string
	Age        int
	Coursename string
	Price      int
}

func main() {
	connection := "user=postgres password=nodirbek dbname=homework sslmode=disable"
	db, err := sql.Open("postgres", connection)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS students(id SERIAL PRIMARY KEY, name VARCHAR, age INT)")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS courses(id SERIAL PRIMARY KEY, courseName VARCHAR, price INT)")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS students_courses(id SERIAL, student_id INT, course_id INT, FOREIGN KEY(student_id) REFERENCES students(ID), FOREIGN KEY (course_id) REFERENCES courses(ID))")
	if err != nil {
		panic(err)
	}

	for {
		var menu int
		fmt.Println("\n1 - Create Student\n2 - Update Student or Course\n3 - Get one student's details\n4 - Get all details\n5 - Delete\n6 - Exit")
		fmt.Scan(&menu)
		switch menu {
		case 1:
			createStudentCourse(db)

		case 2:
			updateStudentCourse(db)

		case 3:
			getOne(db)

		case 4:
			getAll(db)

		case 5:
			delete(db)

		case 6:
			return

		default:
			fmt.Println("Please select appropriate number")
		}
	}
}

func createStudentCourse(db *sql.DB) {
	var studentName string
	var studentAge int
	fmt.Print("Enter student's name: ")
	fmt.Scan(&studentName)
	fmt.Print("Enter student's age: ")
	fmt.Scan(&studentAge)

	var respStudentId int
	var respCourseId int

	db.QueryRow("INSERT INTO students(name, age) VALUES($1, $2) RETURNING id", studentName, studentAge).Scan(&respStudentId)

	var courseName string
	var price int
	fmt.Print("Enter course name: ")
	fmt.Scan(&courseName)
	fmt.Print("Enter price: ")
	fmt.Scan(&price)

	db.QueryRow("INSERT INTO courses(courseName, price) VALUES($1, $2) RETURNING id", courseName, price).Scan(&respCourseId)

	_, err := db.Exec("INSERT INTO students_courses(student_id, course_id) VALUES($1, $2)", respStudentId, respCourseId)
	if err != nil {
		panic(err)
	}
	pp.Print("Succesfully created!!!")
}

func updateStudentCourse(db *sql.DB) {
	var menu int
	fmt.Println("1 - Change student's data\n2 - Change student's course data")
	fmt.Scan(&menu)
	if menu == 1 {
		var studentId int
		var studentName string
		var age int
		fmt.Print("Enter student's ID to update: ")
		fmt.Scan(&studentId)
		fmt.Print("Enter student's name : ")
		fmt.Scan(&studentName)
		fmt.Print("Enter student's age: ")
		fmt.Scan(&age)
		_, err := db.Exec("UPDATE students SET name = $1, age = $2 WHERE id = $3", studentName, age, studentId)
		if err != nil {
			panic(err)
		}
		pp.Print("Succesfully changed!!!")
	} else if menu == 2 {
		var studentId int
		var courseName string
		var price int
		var studentName string

		fmt.Print("Enter student's ID to update: ")
		fmt.Scan(&studentId)

		db.QueryRow("SELECT name FROM students WHERE id = $1", studentId).Scan(&studentName)
		fmt.Printf("You are changing %s's course details!!!", studentName)

		fmt.Print("Entera  new coursename: ")
		fmt.Scan(&courseName)
		fmt.Print("Enter a new price: ")
		fmt.Scan(&price)

		_, err := db.Exec("UPDATE courses SET coursename = $1, price = $2 WHERE id = $3", courseName, price, studentId)
		if err != nil {
			panic(err)
		}
		pp.Print("Succesfully changed!!!")
	} else {
		fmt.Println("Invalid option, please choose 1 or 2 !!!")
	}
}

func getOne(db *sql.DB) {
	var studentId int
	fmt.Print("Enter student's ID to get: ")
	fmt.Scan(&studentId)

	var student Student
	err := db.QueryRow("SELECT * from students WHERE id = $1", studentId).Scan(&student.Id, &student.Name, &student.Age)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ID: %d\tName: %s\tAge: %d", student.Id, student.Name, student.Age)
}

func getAll(db *sql.DB) {
	var data []ReceivedData

	rows, err := db.Query("select s.id, s.name, s.age, c.coursename, c.price from students s join students_courses sc on s.id = sc.student_id join courses c on c.id = course_id")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var rd ReceivedData
		if err := rows.Scan(&rd.Id, &rd.Name, &rd.Age, &rd.Coursename, &rd.Price); err != nil {
			panic(err)
		}
		data = append(data, rd)
	}

	pp.Print(data)
}

func delete(db *sql.DB) {
	var menu int
	fmt.Println("1 - Delete Student\n2 - Delete Course")
	fmt.Print("Enter your choice: ")
	fmt.Scan(&menu)

	switch menu {
	case 1:
		deleteStudent(db)
	case 2:
		deleteCourse(db)
	default:
		fmt.Println("Invalid option, please choose 1 or 2 !!!")
	}
}

func deleteStudent(db *sql.DB) {
	var studentId int
	fmt.Print("Enter student's ID to delete: ")
	fmt.Scan(&studentId)

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM students WHERE id = $1", studentId).Scan(&count)
	if err != nil {
		panic(err)
	}

	if count == 0 {
		fmt.Println("Student not found.")
		return
	}

	_, err = db.Exec("DELETE FROM students_courses WHERE student_id = $1", studentId)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DELETE FROM students WHERE id = $1", studentId)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully deleted the student and associated course details.")
}

func deleteCourse(db *sql.DB) {
	var courseId int
	fmt.Print("Enter course ID to delete: ")
	fmt.Scan(&courseId)

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM courses WHERE id = $1", courseId).Scan(&count)
	if err != nil {
		panic(err)
	}

	if count == 0 {
		fmt.Println("Course not found.")
		return
	}

	_, err = db.Exec("DELETE FROM students_courses WHERE course_id = $1", courseId)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DELETE FROM courses WHERE id = $1", courseId)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully deleted the course and associated student details.")
}
