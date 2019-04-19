package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"image/color"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/ehardi19/go-fuzzy-logic"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func main() {

	// READ THE CSV DATA
	csvFile, _ := os.Open("DataTugas3.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	defer csvFile.Close()
	var data []fuzzy.Number
	reader.Read()
	for {
		var fn fuzzy.Number
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fn.Interview.ID = row[0]
		fn.Interview.Competence, err = strconv.ParseFloat(row[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		fn.Interview.Personality, err = strconv.ParseFloat(row[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		data = append(data, fn)
	}

	// CREATING PREDICTION FILE
	file, _ := os.Create("TebakanTugas3.csv")
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// FUZZY
	employeeAcceptance := fuzzy.EmployeeAcceptance{}
	for i := range data {
		employeeAcceptance.Fuzzification(&data[i])
		employeeAcceptance.Inference(&data[i])
		employeeAcceptance.Defuzzification(&data[i])
	}

	// INSERTING DATA TO PREDICTION
	head := []string{
		"ID",
		"Tes Kompetensi",
		"Kepribadian",
		"Diterima",
	}
	if err := writer.Write(head); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
	for i := range data {
		csvData := []string{
			fmt.Sprintf("%s", data[i].Interview.ID),
			fmt.Sprintf("%.1f", data[i].Interview.Competence),
			fmt.Sprintf("%.1f", data[i].Interview.Personality),
			fmt.Sprintf("%s", data[i].Inference),
		}
		if err := writer.Write(csvData); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// PLOTTING
	scatterData := make(plotter.XYZs, len(data))
	for i := range data {
		scatterData[i].X = data[i].Interview.Competence
		scatterData[i].Y = data[i].Interview.Personality
		scatterData[i].Z = data[i].CrispValue
	}

	minZ, maxZ := math.Inf(1), math.Inf(-1)
	for _, xyz := range scatterData {
		if xyz.Z > maxZ {
			maxZ = xyz.Z
		}
		if xyz.Z < minZ {
			minZ = xyz.Z
		}
	}

	plt, _ := plot.New()
	plt.Title.Text = "Penerimaan Pegawai"
	plt.X.Label.Text = "Tes Kompetensi"
	plt.Y.Label.Text = "Kepribadian"
	plt.Legend.Add("The bigger the circle, the bigger the value")

	sc, err := plotter.NewScatter(scatterData)
	if err != nil {
		panic(err)
	}

	sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		x, y, z := scatterData.XYZ(i)
		c := color.RGBA{R: 190 + uint8(x*20), G: 128 + uint8(y*20), B: 0 + uint8(z*20), A: 255}
		var minRadius, maxRadius = vg.Points(-10), vg.Points(10)
		rng := maxRadius - minRadius
		d := (z - minZ) / (maxZ - minZ)
		r := vg.Length(d) * rng * 2
		return draw.GlyphStyle{Color: c, Radius: r, Shape: draw.CircleGlyph{}}
	}
	plt.Add(sc)

	plt.Save(15*vg.Inch, 15*vg.Inch, "contour.svg")

	// PLOT TES KOMPENTENSI FUNCTION
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Tes Kompetensi Membership Functions"
	p.X.Label.Text = "Tes Kompetensi"
	p.Y.Label.Text = "Fuzzy"

	low := plotter.NewFunction(employeeAcceptance.CompetenceLow)
	low.Color = color.RGBA{R: 255, A: 255}

	mid := plotter.NewFunction(employeeAcceptance.ComptenceMiddle)
	mid.Color = color.RGBA{G: 255, A: 255}

	high := plotter.NewFunction(employeeAcceptance.CompetenceHigh)
	high.Color = color.RGBA{B: 255, A: 255}

	p.Add(low, mid, high)

	p.X.Min = 0
	p.X.Max = 100
	p.Y.Min = 0
	p.Y.Max = 1

	if err := p.Save(5*vg.Inch, 5*vg.Inch, "competenceFunction.svg"); err != nil {
		panic(err)
	}

	// PLOT PERSONALITY FUNCTION
	p, err = plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Kepribadian Membership Functions"
	p.X.Label.Text = "Kepribadian"
	p.Y.Label.Text = "Fuzzy"

	low = plotter.NewFunction(employeeAcceptance.PersonalityLow)
	low.Color = color.RGBA{R: 255, A: 255}

	mid = plotter.NewFunction(employeeAcceptance.PersonalityMiddle)
	mid.Color = color.RGBA{G: 255, A: 255}

	high = plotter.NewFunction(employeeAcceptance.PersonalityHigh)
	high.Color = color.RGBA{B: 255, A: 255}

	p.Add(low, mid, high)

	p.X.Min = 0
	p.X.Max = 100
	p.Y.Min = 0
	p.Y.Max = 1

	// Save the plot to a PNG file.
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "personalityFunction.svg"); err != nil {
		panic(err)
	}
}
