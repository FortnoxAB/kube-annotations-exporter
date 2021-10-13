#!/usr/bin/env groovy
// vim: ft=Jenkinsfile
library 'whatsout'

node('go1.17') {
	container('run'){
		def newTag = ''
		def tag = ''
		def gitTag = ''

		try {
			stage('Checkout'){
				checkout scm
				notifyBitbucket()
				gitTag = sh(script: 'git tag -l --contains HEAD', returnStdout: true).trim()
			}


			stage('Fetch dependencies'){
				// using ID because: https://issues.jenkins-ci.org/browse/JENKINS-32101
				sshagent(credentials: ['18270936-0906-4c40-a90e-bcf6661f501d']) {
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

					sshagent(credentials: ['18270936-0906-4c40-a90e-bcf6661f501d']) {
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
				strippedTag = tag.replaceFirst('v', '')
				stage('Build the application'){
					echo "Building with docker tag ${strippedTag}"
					sh('CGO_ENABLED=0 GOOS=linux go build')
				}

				stage('Generate docker image'){
					image = docker.build('fortnox/kube-annotations-exporter:'+strippedTag, '--pull .')
				}

				stage('Push docker image'){
					docker.withRegistry("https://quay.io", 'docker-registry') {
						image.push()
					}
				}
			}

			currentBuild.result = 'SUCCESS'
		} catch(err) {
			currentBuild.result = 'FAILED'
			notifyBitbucket()
			if (newTag != '') {
				sshagent(credentials: ['18270936-0906-4c40-a90e-bcf6661f501d']) {
					sh("git tag -d ${tag}")
					sh("git push --delete origin ${tag}")
				}
			}
			throw err
		}

		notifyBitbucket()
	}
}

