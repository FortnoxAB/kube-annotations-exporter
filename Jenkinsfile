#!/usr/bin/env groovy
// vim: ft=Jenkinsfile
library 'whatsout'

node('go1.21') {
	container('run'){
		def newTag = ''
		def tag = ''
		def gitTag = ''

		try {
			stage('Checkout'){
				checkout scm
				gitTag = sh(script: 'git tag -l --contains HEAD', returnStdout: true).trim()
			}

			stage('Fetch dependencies'){
				// using ID because: https://issues.jenkins-ci.org/browse/JENKINS-32101
				sshagent(credentials: ['github-key']) {
					//sh('go get -d -v -t ./...')
					sh('go mod download')
				}
			}
			stage('Run test'){
				sh('make test')
			}

			if(gitTag != ''){
				tag = gitTag
			}else if (env.BRANCH_NAME == 'master'){
				stage('create tag') {
					CURR_TAG = sh(
							script: "git for-each-ref --count=1 --sort=-v:refname --format '%(tag)' refs/tags/v*",
							returnStdout: true
					).trim()
					echo "current tag: "+CURR_TAG
					if(CURR_TAG == ''){
						CURR_TAG = 'v0'
					}

					PARTS = CURR_TAG.replaceFirst('v','').split('\\.')
					PARTS[PARTS.length-1] = (PARTS[PARTS.length-1] as int)+1
					newTag = 'v'+PARTS.join(".")

					sshagent(credentials: ['github-key']) {
						sh('git config user.email "fortnoxdrift@fortnox.se"')
						sh('git config user.name "jenkins"')
						sh("git tag -a ${newTag} -m ${newTag}")
						sh("git push origin ${newTag}")
					}
					tag = newTag
				}

				createRelease(CURR_TAG, newTag, 'kube-annotations-exporter', 'kube-annotations-exporter')
			}

			if( tag != ''){
				stage('Build the application'){
					echo "Building with docker tag ${tag}"
					docker.withRegistry("https://quay.io", 'docker-registry') {
						sh("VERSION=${tag} make push")
					}
				}
			}

			currentBuild.result = 'SUCCESS'
		} catch(err) {
			currentBuild.result = 'FAILED'
			if (newTag != '') {
				sshagent(credentials: ['github-key']) {
					sh("git tag -d ${tag}")
					sh("git push --delete origin ${tag}")
				}
			}
			throw err
		}
	}
}

