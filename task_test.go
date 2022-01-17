package gotask

import (
	"container/list"
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		max  int
		done func(j *Job)
	}
	tests := []struct {
		name string
		args args
		want Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Load(tt.args.max, tt.args.done); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_task_AddRunner(t1 *testing.T) {
	type fields struct {
		max      int
		lock     *sync.RWMutex
		idle     *sync.Pool
		jobs     map[interface{}]*list.Element
		ll       *list.List
		doneHook func(job *Job)
	}
	type args struct {
		runner Runner
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:   "",
			fields: fields{},
			args: args{
				runner: &runner{
					name: "runner1",
				},
			},
			wantErr: false,
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				runner: &runner{
					name: "runner2",
				},
			},
			wantErr: false,
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				runner: &runner{
					name: "runner3",
				},
			},
			wantErr: false,
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				runner: &runner{
					name: "runner4",
				},
			},
			wantErr: true,
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				runner: &runner{
					name: "runner5",
				},
			},
			wantErr: true,
		},
	}
	t := Load(3, func(j *Job) {
		fmt.Println("done", j)
	})

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if _, err := t.AddRunner(tt.args.runner); (err != nil) != tt.wantErr {
				t1.Errorf("AddRunner() error = %v, wantErr %v", err, tt.wantErr)
				t.StopJob(tt.args.runner.Key())
			}
			fmt.Println(t.Runs())
		})
	}
}

type runner struct {
	KeyUUID
	name string
}

func (r runner) Run(ctx context.Context, job *Job) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fmt.Println(r.name, job)
		}
	}
}
