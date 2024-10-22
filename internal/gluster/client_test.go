package gluster

import (
	"slices"
	"testing"
)

func TestContainsVolume(t *testing.T) {
	example := "doge"
	testSlice := []string{"wow", "such", example}
	if !slices.Contains(testSlice, example) {
		t.Fatalf("Hasn't found %v in slice %v", example, testSlice)
	}
}

type testCases struct {
	mountOutput string
	expected    []string
}

func TestParseMountOutput(t *testing.T) {
	var tests = []testCases{
		{
			mountOutput: "/dev/mapper/cryptroot on / type ext4 (rw,relatime,data=ordered) \n" +
				"/dev/mapper/cryptroot on /var/lib/docker/devicemapper type ext4 (rw,relatime,data=ordered)",
			expected: []string{"/", "/var/lib/docker/devicemapper"},
		},
		{
			mountOutput: "/dev/mapper/cryptroot on / type ext4 (rw,relatime,data=ordered) \n" +
				"",
			expected: []string{"/"},
		},
	}
	for _, c := range tests {
		mounts, err := ParseMountOutput(c.mountOutput)
		if err != nil {
			t.Error(err)
		}

		for i, mount := range mounts {
			if mount.MountPoint != c.expected[i] {
				t.Errorf("mountpoint is %v and %v was expected", mount.MountPoint, c.expected[i])
			}
		}
	}

}
