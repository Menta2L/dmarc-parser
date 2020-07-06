/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/jedisct1/dlog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/menta2l/dmarc-parser/internal/dmarc"
	ilog "github.com/menta2l/dmarc-parser/internal/log"
	"github.com/menta2l/dmarc-parser/internal/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var connstr string

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var dbhost string = "unix(/var/run/mysqld/mysqld.sock)"
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			dlog.SetLogLevel(dlog.SeverityDebug)
		}
		dlog.Debug("DEBUG: Starging dmarc-parser")
		if viper.IsSet("DBHOST") {
			//
			// unix(/var/run/mysqld/mysqld.sock)
			dbhost = "tcp(" + viper.GetString("DBHOST") + ")"
		}
		if connstr == "" {
			connstr = fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local",
				viper.GetString("DBUSER"), viper.GetString("DBPASS"), dbhost, viper.GetString("DBNAME"))
		}
		dlog.Debug(fmt.Sprintf("DEBUG: connstr = %s", connstr))
		db, err := gorm.Open("mysql", connstr)
		if err != nil {
			dlog.Errorf("Fatal error = %s", err)

		}
		if debug {
			db.LogMode(true)
		} else {
			db.LogMode(false)
		}
		db.SetLogger(&ilog.GormLogger{})
		defer db.Close()
		db.AutoMigrate(&types.DmarcReport{}, &types.DmarcPOReason{}, types.DmarcSPFAuthResult{}, types.DmarcDKIMAuthResult{}, types.DmarcRecord{})
		db.Model(&types.DmarcPOReason{}).AddForeignKey("record_id", "dmarc_records(id)", "CASCADE", "CASCADE")
		db.Model(&types.DmarcSPFAuthResult{}).AddForeignKey("record_id", "dmarc_records(id)", "CASCADE", "CASCADE")
		db.Model(&types.DmarcDKIMAuthResult{}).AddForeignKey("record_id", "dmarc_records(id)", "CASCADE", "CASCADE")
		db.Model(&types.DmarcRecord{}).AddForeignKey("report_id", "dmarc_reports(id)", "CASCADE", "CASCADE")

		err = dmarc.Parse(input, db)
		if err != nil {
			dlog.Errorf("Fatal error = %s", err)
		}
	},
}
var input string

func init() {
	parseCmd.PersistentFlags().StringVar(&input, "input", "stdin", "File path or stdin default: stdin")
	parseCmd.PersistentFlags().StringVar(&connstr, "connstr", "", "Database connection string")

	viper.BindPFlag("input", parseCmd.PersistentFlags().Lookup("input"))
	rootCmd.AddCommand(parseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// parseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
