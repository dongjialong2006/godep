pipeline {
    agent any

    stages {
        stage("Prepare") {
            steps {
                echo 'Prepare..'
            }
        }
        stage("Test") {
            steps {
                echo 'Test..'
            }
        }
        stage("Analysis") {
            steps {
                echo 'Analysis..'
            }
        }
        stage('Build') {
            steps {
                echo "Build.."
            }
        }
        stage('Deploy') {
            parallel {
                stage('OPS') {
                    agent any
                    steps {
                        echo "On Branch A"
                    }
                }
                stage('SMAC') {
                    agent any
                    steps {
                        sh 'make'
                        archiveArtifacts artifacts: 'bin/*', fingerprint: true
                    }
                }
            }
        }
    }
}