pipeline {
    agent any

    stages {
        stage("Prepare") {
            steps {
                echo 'Prepare..'
            }
        }
        stage('Build') {
            parallel {
                stage('BranchA') {
                    agent any
                    steps {
                        echo "On Branch A"
                    }
                }
                stage('BranchB') {
                    agent any
                    steps {
                        sh 'make'
                        archiveArtifacts artifacts: 'bin/*', fingerprint: true
                    }
                }
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
            }
        }
    }
}