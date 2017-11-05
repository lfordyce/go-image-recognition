package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type ClassifyResult struct {
	Filename string        `json:"filename"`
	Labels   []LabelResult `json:"labels"`
}

type LabelResult struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}

var (
	graph  *tf.Graph
	labels []string
)

func main() {
	if err := loadModel(); err != nil {
		log.Fatal(err)
		return
	}

	r := httprouter.New()
	r.GET("/api", testHandler)
	r.POST("/recognize", recognizeHandler)
	r.POST("/series", recognizeHandlerMultipart)
	// log.Fatal(http.ListenAndServe(":8080", r))
	log.Fatal(http.ListenAndServe(":8080", &Server{r}))
}

type Server struct {
	r *httprouter.Router
}

// Server Wrapper with Access-Control-Allowed-Headers
func (s *Server) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	s.r.ServeHTTP(w, r)
}

func loadModel() error {
	// Load inception model
	model, err := ioutil.ReadFile("/model/tensorflow_inception_graph.pb")

	// New Model
	// model, err := ioutil.ReadFile("/model/inception_v3_2016_08_28_frozen.pb")
	if err != nil {
		return err
	}
	graph = tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		return err
	}
	// Load labels
	labelsFile, err := os.Open("/model/imagenet_comp_graph_label_strings.txt")

	// TF-slim labels
	// labelsFile, err := os.Open("/model/imagenet_slim_labels.txt")
	if err != nil {
		return err
	}
	defer labelsFile.Close()
	scanner := bufio.NewScanner(labelsFile)
	// Labels are separated by newlines
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func testHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	responseJSON(w, "Test from GO Api")
}

func recognizeHandlerMultipart(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		responseError(w, "no multipart form", http.StatusBadRequest)
		return
	}

	formdata := r.MultipartForm

	// Get the *fileheaders
	files := formdata.File["images"]

	for i, _ := range files {
		// For each fileheader, get a handle to the actual file
		imageFile, err := files[i].Open()
		defer imageFile.Close()
		if err != nil {
			responseError(w, "Could not open fileheaders", http.StatusBadRequest)
			return
		}

		imageName := strings.Split(files[i].Filename, ".")
		// Return best labels

		var imageBuffer bytes.Buffer
		// Copy image data to a buffer
		io.Copy(&imageBuffer, imageFile)

		// ...
		// Make tensor
		tensor, err := makeTensorFromImage(&imageBuffer, imageName[:1][0])
		if err != nil {
			responseError(w, "Invalid image", http.StatusBadRequest)
			return
		}

		// Run inference
		session, err := tf.NewSession(graph, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		output, err := session.Run(
			map[tf.Output]*tf.Tensor{
				graph.Operation("input").Output(0): tensor,
			},
			[]tf.Output{
				graph.Operation("output").Output(0),
			},
			nil)
		if err != nil {
			responseError(w, "Could not run inference", http.StatusInternalServerError)
			return
		}

		// Return best labels
		responseJSON(w, ClassifyResult{
			Filename: files[i].Filename,
			Labels:   findBestLabels(output[0].Value().([][]float32)[0]),
		})
	}
}

func recognizeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Read image
	imageFile, header, err := r.FormFile("image")
	// Will contain filename and extension
	imageName := strings.Split(header.Filename, ".")
	if err != nil {
		responseError(w, "Could not read image", http.StatusBadRequest)
		return
	}
	defer imageFile.Close()
	var imageBuffer bytes.Buffer
	// Copy image data to a buffer
	io.Copy(&imageBuffer, imageFile)

	// ...
	// Make tensor
	tensor, err := makeTensorFromImage(&imageBuffer, imageName[:1][0])
	if err != nil {
		responseError(w, "Invalid image", http.StatusBadRequest)
		return
	}

	// Run inference
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			graph.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			graph.Operation("output").Output(0),
		},
		nil)
	if err != nil {
		responseError(w, "Could not run inference", http.StatusInternalServerError)
		return
	}

	// Return best labels
	responseJSON(w, ClassifyResult{
		Filename: header.Filename,
		Labels:   findBestLabels(output[0].Value().([][]float32)[0]),
	})
}

type ByProbability []LabelResult

func (a ByProbability) Len() int           { return len(a) }
func (a ByProbability) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByProbability) Less(i, j int) bool { return a[i].Probability > a[j].Probability }

func findBestLabels(probabilities []float32) []LabelResult {
	// Make a list of label/probability pairs
	var resultLabels []LabelResult
	for i, p := range probabilities {
		if i >= len(labels) {
			break
		}
		resultLabels = append(resultLabels, LabelResult{Label: labels[i], Probability: p})
	}
	// Sort by probability
	sort.Sort(ByProbability(resultLabels))
	// Return top 5 labels
	return resultLabels[:5]
}
