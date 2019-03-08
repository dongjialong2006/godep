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
                    steps {
                        echo "On Branch A"
                    }
                }
                stage('BranchB') {
                    steps {
                        echo "On Branch B"
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