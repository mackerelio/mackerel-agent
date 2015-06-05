// +build linux darwin freebsd

package command

func TestRunOnce(t *testing.T) {
	conf := &config.Config{
		Plugin: map[string]config.PluginConfigs{
			"metrics": map[string]config.PluginConfig{
				"metric1": config.PluginConfig{
					Command: "echo test\t1\t1",
				},
			},
			"checks": map[string]config.PluginConfig{
				"check1": config.PluginConfig{
					Command: "echo 1",
				},
			},
		},
	}
	err := RunOnce(conf)
	if err != nil {
		t.Errorf("RunOnce() should be nomal exit: %s", err)
	}
}
