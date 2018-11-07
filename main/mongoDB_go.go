package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type data struct {
	Timestamp string
	Tx string
	Ty string
	Tz string
	Qx string
	Qy string
	Qz string
	Qw string
}

func Compute_performance(count int, average int){
	session,err:=mgo.Dial("")
	if err != nil {panic(err) }
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	col_go := session.DB("test").C("collection_go")

	var timeList_insert []float64
	var timeList_updateOneByOne []float64
	var timeList_updateOneRandom []float64
	var timeList_findOneByOne []float64
	var timeList_findOneRandom []float64
	var timeList_delete []float64

    for i:=0; i<average; i++ {
    	fmt.Println((i+1)," -------------------------------------")
    	//------------- operations ---------------------------
		InsertOne(col_go, count, &timeList_insert)
		UpdateOneByOne(col_go, count, &timeList_updateOneByOne)
		UpdateOne_random(col_go, count, &timeList_updateOneRandom)
		FindOneByOne(col_go, count, &timeList_findOneByOne)
		FindOne_random(col_go, count, &timeList_findOneRandom)
		DeleteOne(col_go, count, &timeList_delete)

		// After testing, delete the collection from database!!!
		col_go.DropCollection()
	}

	// compute the average time consuming by different operations
	averageTime_insert := Compute_average_value(timeList_insert)
	averageTime_updateOneByOne := Compute_average_value(timeList_updateOneByOne)
	averageTime_updateOneRandom := Compute_average_value(timeList_updateOneRandom)
	averageTime_findOneByOne := Compute_average_value(timeList_findOneByOne)
	averageTime_findRandom := Compute_average_value(timeList_findOneRandom)
	averageTime_delete := Compute_average_value(timeList_delete)

	fmt.Println("======= " + strconv.Itoa(average) + " times Computation ===========")
	fmt.Printf("averageTime_insert: %f[s]\n", averageTime_insert)
	fmt.Printf("averageTime_updateOneByOne: %f[s]\n", averageTime_updateOneByOne)
	fmt.Printf("averageTime_updateOneRandom: %f[s]\n", averageTime_updateOneRandom)
	fmt.Printf("averageTime_findOneByOne: %f[s]\n", averageTime_findOneByOne)
	fmt.Printf("averageTime_findRandom: %f[s]\n", averageTime_findRandom)
	fmt.Printf("averageTime_delete: %f[s]\n", averageTime_delete)

	// write the results to file
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	out_path := dir + "/results/"+ fmt.Sprintf("%d",count) +"_data_timeConsuming.txt" //TODO
	Write_to_txt(out_path, averageTime_insert, averageTime_updateOneByOne,
		         averageTime_updateOneRandom, averageTime_findOneByOne,
		         averageTime_findRandom, averageTime_delete)
}

func InsertOne(col_go *mgo.Collection, count int, timeList_insert *[]float64){
	start := time.Now()

	// Insert data one by one.
	for num_insert :=0; num_insert<count; num_insert++ {
		d := data{
			Timestamp: strconv.Itoa(num_insert),
			Tx: "-0.00105005",
			Ty: "-0.00159259",
			Tz: "0",
			Qx: "0",
			Qy: "0",
			Qz: "0.00128360443",
			Qw: "0.999999176",
		}
		err := col_go.Insert(&d)
		if err != nil {
			log.Fatal(err)
		}
	}

	end := time.Now()
	duration := end.Sub(start).Seconds()
	fmt.Println("Insert operation time consuming: ", duration)
	*timeList_insert = append(*timeList_insert, duration)
}

func UpdateOneByOne(col_go *mgo.Collection, count int, timeList_updateOneByOne *[]float64){
	start := time.Now()
	for num_update :=0; num_update<count; num_update++ {
		err := col_go.Update(bson.M{"timestamp": strconv.Itoa(num_update)},
		       bson.M{"$set": bson.M{"qx": "1"}})

		if err != nil{
			log.Fatal(err)
		}
	}

	end := time.Now()
	duration := end.Sub(start).Seconds()
	fmt.Println("Update operation(one by one) time consuming: ", duration)
	*timeList_updateOneByOne = append(*timeList_updateOneByOne, duration)
}

func UpdateOne_random(col_go *mgo.Collection, count int, timeList_updateOneRandom *[]float64){
	start := time.Now()

	for num_update :=0; num_update<count; num_update++ {
		r := rand.Intn(count)

		err := col_go.Update(bson.M{"timestamp": strconv.Itoa(r)},
			bson.M{"$set": bson.M{"qx": "2"}})

		if err != nil{
			log.Fatal(err)
		}
	}

	end := time.Now()
	duration := end.Sub(start).Seconds()
	fmt.Println("Update operation(randomly) time consuming: ", duration)
	*timeList_updateOneRandom = append(*timeList_updateOneRandom, duration)
}

func FindOneByOne(col_go *mgo.Collection, count int, timeList_findOneByOne *[]float64){
	start := time.Now()
    tmp := data{}

	for num_update :=0; num_update<count; num_update++ {
		err := col_go.Find(bson.M{"timestamp": strconv.Itoa(num_update)}).One(&tmp)
		if err != nil{
			log.Fatal(err)
		}
	}

	end := time.Now()
	duration := end.Sub(start).Seconds()
	fmt.Println("Find operation(one by one) time consuming: ", duration)
	*timeList_findOneByOne = append(*timeList_findOneByOne, duration)
}

func FindOne_random(col_go *mgo.Collection, count int, timeList_findOneRandom *[]float64){
	start := time.Now()
	tmp := data{}
    r := rand.Intn(count)

	for num_update :=0; num_update<count; num_update++ {
		err := col_go.Find(bson.M{"timestamp": strconv.Itoa(r)}).One(&tmp)

		if err != nil{
			log.Fatal(err)
		}
	}

	end := time.Now()
	duration := end.Sub(start).Seconds()
	fmt.Println("Find operation(randomly) time consuming: ", duration)
	*timeList_findOneRandom = append(*timeList_findOneRandom, duration)
}

func DeleteOne(col_go *mgo.Collection, count int, timeList_delete *[]float64){
	start := time.Now()

	for num_update :=0; num_update<count; num_update++ {
		err := col_go.Remove(bson.M{"timestamp": strconv.Itoa(num_update)})

		if err != nil{
			log.Fatal(err)
		}
	}

	end := time.Now()
	duration := end.Sub(start).Seconds()
	fmt.Println("Delete operation time consuming: ", duration)
	*timeList_delete = append(*timeList_delete, duration)
}

func Compute_average_value(list []float64) float64{
	var sum float64 = 0.0
	for _, value:= range list{
		sum += value
	}
	l := float64(len(list))
	return sum / l
}

func Write_to_txt(out_path string, averageTime_insert float64,
	              averageTime_updateOneByOne float64,
	              averageTime_updateOneRandom float64,
	              averageTime_findOneByOne float64,
	              averageTime_findRandom float64,
	              averageTime_delete float64) {
	f, err := os.Create(out_path)

	if err != nil{
		panic(err)
	}

    defer f.Close()

    f.WriteString("averageTime_insert: ")
    f.WriteString(fmt.Sprintf("%f", averageTime_insert))
    f.WriteString("[s]\n")

	f.WriteString("averageTime_updateOneByOne: ")
	f.WriteString(fmt.Sprintf("%f", averageTime_updateOneByOne))
	f.WriteString("[s]\n")

	f.WriteString("averageTime_updateOneRandom: ")
	f.WriteString(fmt.Sprintf("%f", averageTime_updateOneRandom))
	f.WriteString("[s]\n")

	f.WriteString("averageTime_findOneByOne: ")
	f.WriteString(fmt.Sprintf("%f", averageTime_findOneByOne))
	f.WriteString("[s]\n")

	f.WriteString("averageTime_findRandom: ")
	f.WriteString(fmt.Sprintf("%f", averageTime_findRandom))
	f.WriteString("[s]\n")

	f.WriteString("averageTime_delete: ")
	f.WriteString(fmt.Sprintf("%f", averageTime_delete))
	f.WriteString("[s]\n")
}