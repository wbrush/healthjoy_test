//  +build integration

package integrations

import (
	"fmt"
	"github.com/wbrush/go-template-service/dao/postgres"
	"github.com/wbrush/go-template-service/models"
	"net/url"
	"time"
)

/*
    It seems like a good idea to create a different shard for each test type. This will allow us to review the data once
    the test completes without having to worry about other tests corrupting or changing the data we are interested in.
	Please try to continue this going forward.
*/
const CreateTemplateShard = int64(1) //this shard is used for record creation test, so after init it must be empty
const GetTemplateShard = int64(2)
const ListTemplateShard = int64(3)
const DeleteTemplateShard = int64(4)
const UpdateTemplateShard = int64(5)

const testRecordNum = int64(1)

func MakeHeaders(shard int64) map[string]string {

	sh := datamodels.X_OptiiUserInfo{
		Shard: make([]int64, 0),
	}
	sh.Shard = append(sh.Shard, shard)

	//  set up headers
	headers := make(map[string]string)
	headers["X-Optii-UserInfo"], _ = sh.ToString()

	return headers
}

func InitializeDBs() {
	var shardMap = make(map[int64]string)
	var shardNames = []string{"Test record creation",
		"Test record get by id",
		"Test record get by list",
		"Test record delete by id",
		"Test record update by id",
	}

	//  initialize the DBs here; would prefer to send the pubsub messages to create shards but don't have pubsub working at the moment.
	fmt.Println("Initializing the shards!")

	dao := postgres.GetPgDao()
	shards, err := dao.GetShardList()
	for i := range shards {
		shardMap[int64(i+1)] = "item"
	}

	if err != nil && err.Error() != db.ErrNoShardsYet {
		fmt.Errorf("error reading shard list: %s", err.Error())
	} else if len(shards) > 0 {
		fmt.Println("Have shards! Verifying we have the right ones")
	} else {
		fmt.Println("No shards! Need to initialize")
	}

	for i := range shardNames {
		shardId := int64(i + 1)
		_, ok := shardMap[shardId]
		if !ok {
			shard := &datamodels.Shard{
				ShardId:      shardId,
				PropertyName: shardNames[i],
				CreatedAt:    time.Time{},
				UpdatedAt:    nil,
				ArchivedAt:   nil,
			}

			err = dao.AddNewShard(shard)
			if err != nil {
				fmt.Errorf("error creating shard %d: %s", shard.ShardId, err.Error())
			}

			//create test templates
			if shardId != 1 { //shard 1 is for create test only
				rec := tdatamodels.Template{
					Id:           testRecordNum,
					Name:         "Test template record",
					Status:       tdatamodels.TemplateStatusNew,
					TemplateSelf: "self",
				}

				_, err = dao.CreateTemplate(shardId, &rec)
				if err != nil {
					fmt.Errorf("error adding back record %d in shard %d: %s", rec.Id, shardId, err.Error())
				}
			}
		} else {
			//  need to be sure that shard is empty
			filters := url.Values{}
			list, _, _, _ := dao.ListTemplates(shardId, filters)

			foundDelete := false
			for _, template := range list {
				if template.Id != testRecordNum {
					foundDelete = true
				}
				_, err = dao.DeleteTemplateById(shardId, template.Id)
				if err != nil {
					fmt.Errorf("error deleting record %d in shard %d: %s", template.Id, shardId, err.Error())
				}
			}

			if !foundDelete {
				fmt.Printf("Creating record on shard %d \n", shardId)
				rec := tdatamodels.Template{
					Id:           testRecordNum,
					Name:         "Replacing missing record",
					Status:       tdatamodels.TemplateStatusNew,
					TemplateSelf: "self",
				}
				_, err = dao.CreateTemplate(shardId, &rec)
				if err != nil {
					fmt.Errorf("error adding back record %d in shard %d: %s", rec.Id, shardId, err.Error())
				}
			}
		}
	}

}
