/*
*
  - Created by GoLand.
  - User: buzzlight.frank@qq.com
  - Date: 2025/3/15
  - Time: 14:12
    This is a task to be exectued in a specify fz;
*/
package task

import (
	"ThinkTankCentral/global"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func RegisterScheduledTasks(c *cron.Cron) error {
	if _, err := c.AddFunc("@hourly", func() {
		if err := UpdateArticleViewsSyncTask(); err != nil {
			global.Log.Error("Failed to update article views:", zap.Error(err))
		}
	}); err != nil {
		return err
	}
	return nil
}
