package plot

import (
	"net/http"

	"github.com/andewx/go-iio/sdr"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

// FFTPlotServer is a struct to hold the FFT data
type FFTPlotServer struct {
	Host    string
	Port    string
	Title   string
	XLabel  string
	YLabel  string
	Data    []complex64
	Domain  func(int) float64
	handler func(w http.ResponseWriter, r *http.Request)
	SDR     *sdr.SDR
}

// NewFFTPlotServer creates a new FFTPlotServer
func NewFFTPlotServer(sdr *sdr.SDR, data []complex64) *FFTPlotServer {
	h := &FFTPlotServer{
		Host:   "localhost",
		Port:   "8080",
		Title:  "FFT Plot",
		XLabel: "Frequency",
		YLabel: "Magnitude",
		SDR:    sdr,
		Data:   data,
		Domain: nil,
	}

	h.Domain = h.callbackDomain

	return h
}

func (svr *FFTPlotServer) callbackDomain(i int) float64 {

	bw := svr.SDR.Params.SampleRate / 2
	cf := svr.SDR.Params.CarrierFrequency
	stp := svr.SDR.Params.SampleRate / float64(svr.SDR.Params.FourierPoint)
	return cf - bw + (float64(i) * stp)
}

func generateFFTData(data []complex64, isReal bool) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(data); i++ {
		if isReal {
			items = append(items, opts.LineData{Value: real(data[i])})
		} else {
			items = append(items, opts.LineData{Value: imag(data[i])})
		}
	}
	return items
}

func generateDomainData(domain func(int) float64, length int) []float64 {
	items := make([]float64, 0)
	for i := 0; i < length; i++ {
		items = append(items, domain(i))
	}
	return items
}

// FFTServer creates a web server to display the FFT data
func (svr *FFTPlotServer) handle(w http.ResponseWriter, r *http.Request) {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeChalk}),
		charts.WithTitleOpts(opts.Title{
			Title:    svr.Title,
			Subtitle: "Complex FFT",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: svr.XLabel,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: svr.YLabel,
		}),
	)

	// Set X Axis Data with Callback Function
	xAxisData := generateDomainData(svr.Domain, len(svr.Data))
	fftRealData := generateFFTData(svr.Data, true)
	fftImagData := generateFFTData(svr.Data, false)

	line.SetXAxis(xAxisData).
		AddSeries("Real", fftRealData).
		AddSeries("Imaginary", fftImagData).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(false)}))
	line.Render(w)
}

// Start starts the FFTPlotServer
func (svr *FFTPlotServer) Start() {
	http.HandleFunc("/", svr.handle)
	http.ListenAndServe(svr.Host+":"+svr.Port, nil)
}
