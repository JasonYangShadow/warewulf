package list

import (
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/warewulf/warewulf/internal/app/wwctl/table"
	apiimage "github.com/warewulf/warewulf/internal/pkg/api/image"
	"github.com/warewulf/warewulf/internal/pkg/image"
	"github.com/warewulf/warewulf/internal/pkg/util"
)

var imageList = apiimage.ImageList

func CobraRunE(vars *variables) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		t := table.New(cmd.OutOrStdout())
		showSize := vars.size || vars.chroot || vars.compressed
		if showSize || vars.full || vars.kernel {
			imageInfo, err := imageList()
			if err != nil {
				return err
			}
			if vars.full {
				t.AddHeader("IMAGE NAME", "NODES", "KERNEL VERSION", "CREATION TIME", "MODIFICATION TIME", "SIZE")
				for i := 0; i < len(imageInfo); i++ {
					createTime := time.Unix(int64(imageInfo[i].CreateDate), 0)
					modTime := time.Unix(int64(imageInfo[i].ModDate), 0)
					sz := util.ByteToString(int64(imageInfo[i].ImgSize))
					if vars.compressed {
						sz = util.ByteToString(int64(imageInfo[i].ImgSizeComp))
					}
					if vars.chroot {
						sz = util.ByteToString(int64(imageInfo[i].Size))
					}
					t.AddLine(
						imageInfo[i].Name,
						strconv.FormatUint(uint64(imageInfo[i].NodeCount), 10),
						imageInfo[i].KernelVersion,
						createTime.Format(time.RFC822),
						modTime.Format(time.RFC822),
						sz,
					)
				}
			} else if vars.kernel {
				t.AddHeader("IMAGE NAME", "NODES", "KERNEL VERSION")
				for i := 0; i < len(imageInfo); i++ {
					t.AddLine(
						imageInfo[i].Name,
						strconv.FormatUint(uint64(imageInfo[i].NodeCount), 10),
						imageInfo[i].KernelVersion,
					)
				}
			} else if showSize {
				t.AddHeader("IMAGE NAME", "NODES", "SIZE")
				for i := 0; i < len(imageInfo); i++ {
					sz := util.ByteToString(int64(imageInfo[i].ImgSize))
					if vars.compressed {
						sz = util.ByteToString(int64(imageInfo[i].ImgSizeComp))
					}
					if vars.chroot {
						sz = util.ByteToString(int64(imageInfo[i].Size))
					}

					t.AddLine(
						imageInfo[i].Name,
						strconv.FormatUint(uint64(imageInfo[i].NodeCount), 10),
						sz,
					)
				}
			}
		} else {
			t.AddHeader("IMAGE NAME")
			list, _ := image.ListSources()
			for _, cont := range list {
				t.AddLine(cont)
			}
		}
		t.Print()
		return
	}
}
