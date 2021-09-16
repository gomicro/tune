package tune

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// CWStatter represents all the elements necessary to meet the Statter struct
// used by the tune statting client
type CWStatter struct {
	prefix string
	svc    *cloudwatch.CloudWatch
}

// NewCWStatter takes a prefix and returns a configured statter to point at
// cloudwatch
func NewCWStatter(prefix, region string) (*Client, error) {
	sess, err := session.NewSession(&aws.Config{Region: &region})
	if err != nil {
		return nil, fmt.Errorf("failed to create new session: %v", err.Error())
	}

	svc := cloudwatch.New(sess)

	c := &Client{
		&CWStatter{
			prefix: prefix,
			svc:    svc,
		},
	}

	return c, nil
}

// NewCWStatterWithClient takes a prefix and http client and returns a
// configured statter to point at cloudwatch
func NewCWStatterWithClient(prefix, region string, client *http.Client) (*Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:     &region,
		HTTPClient: client,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create session for AWS client: %v", err.Error())
	}

	svc := cloudwatch.New(sess)

	c := &Client{
		&CWStatter{
			prefix: prefix,
			svc:    svc,
		},
	}

	return c, nil
}

func metricNames(bucket string) (string, string, string) {
	s := strings.Split(bucket, ".")

	switch {
	case len(s) == 2:
		return s[0], strings.Join(s[1:], "."), strings.Join(s[1:], ".")
	case len(s) >= 3:
		return s[0], strings.Join(s[1:len(s)-1], "."), s[len(s)-1]
	default:
		return s[0], s[0], s[0]
	}
}

// Counter meets the Statter interface
func (cws *CWStatter) Counter(sampleRate float32, bucket string, n ...int) {
	dname, dvalue, mname := metricNames(bucket)

	unit := "Count"
	dimension := []*cloudwatch.Dimension{&cloudwatch.Dimension{
		Name:  &dname,
		Value: &dvalue,
	}}

	datums := make([]*cloudwatch.MetricDatum, 0, len(n))
	for _, v := range n {
		d := &cloudwatch.MetricDatum{
			MetricName:        &mname,
			Unit:              &unit,
			StorageResolution: aws.Int64(int64(sampleRate)),
			Value:             aws.Float64(float64(v)),
			Dimensions:        dimension,
		}

		datums = append(datums, d)
	}

	go cws.svc.PutMetricData(&cloudwatch.PutMetricDataInput{ //nolint,errcheck
		Namespace:  &cws.prefix,
		MetricData: datums,
	})
}

// Timing meets the Statter interface
func (cws *CWStatter) Timing(sampleRate float32, bucket string, d ...time.Duration) {
	dname, dvalue, mname := metricNames(bucket)

	unit := "Milliseconds"
	dimension := []*cloudwatch.Dimension{&cloudwatch.Dimension{
		Name:  &dname,
		Value: &dvalue,
	}}

	datums := make([]*cloudwatch.MetricDatum, 0, len(d))
	for _, t := range d {
		d := &cloudwatch.MetricDatum{
			MetricName:        &mname,
			Unit:              &unit,
			StorageResolution: aws.Int64(int64(sampleRate)),
			Value:             aws.Float64(float64(time.Millisecond * t)),
			Dimensions:        dimension,
		}

		datums = append(datums, d)
	}

	go cws.svc.PutMetricData(&cloudwatch.PutMetricDataInput{ //nolint,errcheck
		Namespace:  &cws.prefix,
		MetricData: datums,
	})
}

// Gauge meets the Statter interface
func (cws *CWStatter) Gauge(sampleRate float32, bucket string, value ...string) {
	dname, dvalue, mname := metricNames(bucket)

	unit := "Count"
	dimension := []*cloudwatch.Dimension{&cloudwatch.Dimension{
		Name:  &dname,
		Value: &dvalue,
	}}

	datums := make([]*cloudwatch.MetricDatum, 0, len(value))
	for _, v := range value {
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			d := &cloudwatch.MetricDatum{
				MetricName:        &mname,
				Unit:              &unit,
				StorageResolution: aws.Int64(int64(sampleRate)),
				Value:             &f,
				Dimensions:        dimension,
			}

			datums = append(datums, d)
		}
	}

	go cws.svc.PutMetricData(&cloudwatch.PutMetricDataInput{ //nolint,errcheck
		Namespace:  &cws.prefix,
		MetricData: datums,
	})
}
