module.exports = {
    "branches": ["main"],
    "plugins": [
        [
            "@semantic-release/commit-analyzer",
            {
                "preset": "angular",
                "parserOpts": {
                    "noteKeywords": [
                        "BREAKING CHANGE",
                        "BREAKING CHANGES",
                        "BREAKING"
                    ]
                },
                "releaseRules": [
                    {
                        "type": "build",
                        "release": "patch"
                    },
                    {
                        "type": "deploy",
                        "release": "patch"
                    },
                    {
                        "type": "chore",
                        "release": "patch"
                    },
                    {
                        "type": "docs",
                        "release": "patch"
                    },
                    {
                        "type": "test",
                        "release": "patch"
                    },
                    {
                        "type": "style",
                        "release": "patch"
                    },
                    {
                        "type": "ci",
                        "release": "patch"
                    }
                ]
            }
        ],
        [
            "@semantic-release/release-notes-generator",
            {
                "preset": "conventionalcommits",
                "parserOpts": {
                    "noteKeywords": [
                        "BREAKING CHANGE",
                        "BREAKING CHANGES",
                        "BREAKING"
                    ]
                },
                "writerOpts": {
                    "commitsSort": [
                        "subject",
                        "scope"
                    ]
                },
                "presetConfig": {
                    "types": [
                        {
                            "type": "feat",
                            "section": "Features"
                        },
                        {
                            "type": "fix",
                            "section": "Bug Fixes"
                        },
                        {
                            "type": "chore",
                            "section": "Chores",
                            "hidden": false
                        },
                        {
                            "type": "refactor",
                            "section": "Internal",
                            "hidden": false
                        },
                        {
                            "type": "perf",
                            "section": "Performance",
                            "hidden": false
                        },
                        {
                            "type": "docs",
                            "section": "Documentation",
                            "hidden": false
                        },
                        {
                            "type": "ci",
                            "section": "DevOps",
                            "hidden": false
                        },
                        {
                            "type": "test",
                            "section": "Tests",
                            "hidden": false
                        }
                    ]
                }
            }
        ],
        "@semantic-release/github"
    ],
    "repositoryUrl": "https://github.com/nilpntr/gluster-exporter.git",
}