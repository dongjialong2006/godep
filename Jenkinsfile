pipeline {
    agent any
    stages {
        stage('Prepare') {
            steps {
                echo 'Prepare..'
            }
        }
        stage('Test') {
            steps {
                echo "Test.."
            }
        }
        stage("Analysis") {
            steps {
                echo 'Analysis..'
            }
        } 
        stage('Dispatch') {
        	parallel {
				stage('OPS') {
                    stages {
		               stage('Init') {
		                   steps {
		                       echo "Init.."
		                   }
		               }
		               stage('Build') {
		                   steps {
		                       echo "Build.."
		                   }
		               }
		               stage('Deploy') {
		                   steps {
		                       echo "Deploy.."
		                   }
		               }
		            }
                }
                stage('ZSY') {
                    stages {
		               stage('Init') {
		                   steps {
		                       echo "Init.."
		                   }
		               }
		               stage('Build') {
		                   steps {
		                       echo "Build.."
		                   }
		               }
		               stage('Deploy') {
		                   steps {
		                       echo "Deploy.."
		                   }
		               }
		            }
                }
                stage('SMAC') {
                    stages {
		               stage('Init') {
		                   steps {
		                       echo "Init.."
		                   }
		               }
		               stage('Build') {
		                   steps {
		                       echo "Build.."
		                   }
		               }
		               stage('Deploy') {
		                   steps {
		                       echo "Deploy.."
		                   }
		               }
		            }
                }
			}
        }
    }
}