package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type IJobThing interface {
	GetJobIds() []int64
	SetJobs(jobs []Job)
}

type JobThing struct {
	Jobs   []Job   `json:"-"`
	JobIds []int64 `json:"job_ids" gorm:"-"`
}

func (jt *JobThing) GetJobIds() []int64 {
	job_ids := make([]int64, len(jt.Jobs))
	for i, j := range jt.Jobs {
		job_ids[i] = j.Id
	}
	return job_ids
}
func (jt *JobThing) SetJobs(jobs []Job) {
	jt.Jobs = jobs
}

func (jt *JobThing) ConvertIdsToJobs(db *gorm.DB) error {
	jobs := make([]Job, len(jt.GetJobIds()))
	for i, jid := range jt.GetJobIds() {
		job := Job{}
		if err := db.First(&job, jid).Error; err != nil {
			return err
		}
		jobs[i] = job
	}
	jt.SetJobs(jobs)
	return nil
}

// func (jt *JobThing) UpdateJobThingFromJson(data map[string]interface{}, db *gorm.DB) error {
// 	value, ok := data["job_ids"]
// 	if ok {
// 		job_ids := []int64{}
// 		for _, num := range value.([]interface{}) { // []interface{}
// 			job_ids = append(job_ids, int64(num.(float64)))
// 		}
// 		jt.JobIds = job_ids
// 		if err := jt.ConvertIdsToJobs(db); err != nil {
// 			return errors.New("invalid associated job_ids")
// 		}
// 	}
// 	// if not "ok", no worries. they weren't updating that association
// 	return nil
// }
