package main

import (
	"fmt"
	"github.com/maltegrosse/go-solpos"
	"strconv"
	"time"

	"os"
	"text/tabwriter"
)

func main() {
	tmpMap := make(map[string]interface{})
	tmpMap["temp"] = 27.0
	tmpMap["press"] = 1006.0
	tmpMap["tilt"] = 33.65
	tmpMap["aspect"] = 135.0
	loc, err := time.LoadLocation("America/Atikokan")
	if err != nil {
		fmt.Println(err)
		return
	}
	dt := time.Date(1999, 7, 22, 9, 45, 37, 0, loc)
	sp, err := solpos.NewSolpos(dt, 33.65, -84.43, tmpMap)
	if err != nil {
		fmt.Println(err)
		return
	}
	nrelMap := make(map[string]interface{})
	resultMap := make(map[string]interface{})
	nrelMap["year"] = "1999"
	resultMap["year"] = sp.GetYear()
	nrelMap["month"] = "07"
	resultMap["month"] = sp.GetMonth()
	nrelMap["daynum"] = "203"
	resultMap["daynum"] = sp.GetDaynum()
	nrelMap["day"] = "22"
	resultMap["day"] = sp.GetDay()
	nrelMap["amass"] = "1.335752"
	resultMap["amass"] = sp.GetAmass()
	nrelMap["ampress"] = "1.326522"
	resultMap["ampress"] = sp.GetAmpress()
	nrelMap["azim"] = "97.032875"
	resultMap["azim"] = sp.GetAzim()
	nrelMap["cosinc"] = "0.912569"
	resultMap["cosinc"] = sp.GetCosinc()
	nrelMap["elevref"] = "48.409931"
	resultMap["elevref"] = sp.GetElevref()
	nrelMap["etr"] = "989.668518"
	resultMap["etr"] = sp.GetEtr()
	nrelMap["etrn"] = "1323.239868"
	resultMap["etrn"] = sp.GetEtrn()
	nrelMap["etrtilt"] = "1207.547363"
	resultMap["etrtilt"] = sp.GetEtrtilt()
	nrelMap["prime"] = "1.037040"
	resultMap["prime"] = sp.GetPrime()
	nrelMap["sbcf"] = "1.201910"
	resultMap["sbcf"] = sp.GetSbcf()
	nrelMap["sunrise"] = "347.173431"
	resultMap["sunrise"] = sp.GetSretr()

	nrelMap["sunset"] = "1181.111206"
	resultMap["sunset"] = sp.GetSsetr()

	nrelMap["unprime"] = "0.964283"
	resultMap["unprime"] = sp.GetUnprime()

	nrelMap["zenref"] = "41.590069"
	resultMap["zenref"] = sp.GetZenref()

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	_, err = fmt.Fprintln(writer, "- \tNREL\tSOLPOS\tDiff")
	for key, _ := range nrelMap {
		a, _ := strconv.ParseFloat(fmt.Sprint(nrelMap[key]), 64)

		b, _ := strconv.ParseFloat(fmt.Sprint(resultMap[key]), 64)
		diff := a - b
		if diff < 0.000001 {
			diff = 0
		}
		_, err = fmt.Fprintln(writer, key, "\t", nrelMap[key], "\t", resultMap[key], "\t", diff)
	}
	err = writer.Flush()

	sp.SetFunction(0)
	sp.SetFunction(solpos.LAmass | solpos.LDoy) // call only the airmass function
	sp.SetPress(1013.0)                         // set your own pressure
	fmt.Println("NREL    -> 37.92  5.59  2.90  1.99  1.55  1.30  1.15  1.06  1.02  1.00")
	zenref := 90.0

	fmt.Print("SOLPOS  -> ")
	for zenref >= 0 {
		sp.SetZenref(zenref)
		err = sp.Calculate()
		if err != nil {
			fmt.Println(err)
			return
		}
		tmpVal := sp.GetAmass()
		fmt.Printf("%.2f", tmpVal)
		fmt.Print("  ")
		zenref -= 10.0

	}
	fmt.Println()
	// 5:45
	fmt.Println(sp.GetSunrise())
	// 19:35
	fmt.Println(sp.GetSunset())

}
