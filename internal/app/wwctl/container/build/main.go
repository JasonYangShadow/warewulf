package build

import (
	"fmt"

	"github.com/spf13/cobra"
	cexec "github.com/warewulf/warewulf/internal/app/wwctl/container/exec"
	"github.com/warewulf/warewulf/internal/pkg/api/container"
	"github.com/warewulf/warewulf/internal/pkg/api/routes/wwapiv1"
	pkgcontianer "github.com/warewulf/warewulf/internal/pkg/container"
	"github.com/warewulf/warewulf/internal/pkg/kernel"
	"github.com/warewulf/warewulf/internal/pkg/util"
)

func CobraRunE(cmd *cobra.Command, args []string) error {
	cbp := &wwapiv1.ContainerBuildParameter{
		ContainerNames: args,
		Force:          BuildForce,
		All:            BuildAll,
		Default:        SetDefault,
		Initramfs:      Initramfs,
	}
	if Initramfs {
		return runInitramfsBuild(cmd, cbp)
	}
	return container.ContainerBuild(cbp)
}

func runInitramfsBuild(cmd *cobra.Command, cbp *wwapiv1.ContainerBuildParameter) (err error) {
	// TODO here we need to bind dracut.module file
	// cexec.SetBinds()
	if cbp == nil {
		return fmt.Errorf("ContainerBuildParameter is nill")
	}

	var containers []string
	if cbp.All {
		containers, err = pkgcontianer.ListSources()
	} else {
		containers = cbp.ContainerNames
	}

	if len(containers) == 0 {
		return
	}

	for _, c := range containers {
		// kernel version, we need to set container kernel version as by default, it'll build against
		// host kernel version, which usually does not exist inside container
		var kver string
		rootfsDir := pkgcontianer.RootFsDir(c)
		kver, err = kernel.FindKernelVersion(rootfsDir)
		if err != nil {
			return fmt.Errorf("failed to locate container kernel version: %s", err)
		}

		err = cexec.CobraRunE(cmd, []string{c, "/usr/bin/dracut --no-hostonly --force --verbose --kver " + kver})
		if err != nil {
			return
		}

		// this will cause the rebuild of container image and then we can retrieve the built initramfs from chroot
		err = util.CopyFile(pkgcontianer.InitramfsBootPath(c, kver), pkgcontianer.InitramfsProvisionPath(kver))
		if err != nil {
			return
		}
	}
	return
}
