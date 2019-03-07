import java.util.Date
import java.text.SimpleDateFormat

node('godep') {
    try{
        echo "prepare"
        stage("Prepare") {
            prepare()
        }
        echo "test"
        stage("Test") {
            // test()
            sh "test success"
        }
        echo "build"
        stage("Build") {
            sh 'make'
        }
        echo "tag"
        stage("Tag"){
            tag("v1.0.0")
        }
    }finally{
        notifyBuild()
    }
}

def prepare(){
    checkout scm
    if(env.BRANCH_NAME.indexOf("/") == -1){
        currentBuild.result = 'ABORTED'
        error('The branches do not need ci.')
    }
}

def notifyBuild() {
    echo currentBuild.currentResult
    def commitUserName = sh(script:"git log -1 --pretty=format:'%an'", returnStdout:true)
    def commitUserEmail = sh(script:"git log -1 --pretty=format:'%ae'", returnStdout:true)
    def commitMessage = sh(script:"git log -1 --pretty=format:'%s'", returnStdout:true)
    def commitID = sh(script:"git log -1 --pretty=format:'%H'", returnStdout:true)
    def commitTime = sh(script:"git log -1 --pretty=format:'%ai'", returnStdout:true)
    def subject = "Jenkins自动构建通知"
    def details = """
                    <html>
                        <body>
                            <div>${JOB_NAME}: ${currentBuild.currentResult}</div>
                            <div>========================= git ====================</div>
                            <div>最近提交者: ${commitUserName}</div>
                            <div>构建分支: ${env.BRANCH_NAME}</div>
                            <div>变更说明: ${commitMessage}</div>
                            <div>commit: ${commitID}</div>
                            <div>提交时间: ${commitTime}</div>
                            <div>==================== Jenkins =====================</div>
                            <div>本次构建结果: ${currentBuild.currentResult}</div>
                            <div><a href='${env.JOB_DISPLAY_URL}'>点击查看构建详情</a></div>
                            <div><a href='${env.JOB_DISPLAY_URL}/tests'>点击查看测试报告</a></div>
                        </body>
                    </html>
                """
    emailext (
        subject: subject,
        body: details,
        mimeType: 'text/html',
        to: commitUserEmail,
        recipientProviders: [[$class: 'DevelopersRecipientProvider']]
    )
}

def tag(String tag){
    sh("git tag | xargs -I {} git tag -d {} || true")
    def GIT_USER = ""
    def GIT_PW = ""
    withCredentials([usernamePassword(credentialsId: 'CI-SNC', passwordVariable: 'GIT_PASSWORD', usernameVariable: 'GIT_USERNAME')]) {
            def u = URLEncoder.encode(GIT_USERNAME,'UTF-8')
            def p = URLEncoder.encode(GIT_PASSWORD,'UTF-8')
            GIT_USER = u
            GIT_PW = p
        }
    def result = 0
    def dateFormat = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss")
    def now = dateFormat.format(new Date())
    result = sh returnStatus: true, script: "git tag ${tag} -m 'auto release ${now}'  &&  git push https://${GIT_USER}:${GIT_PW}@git-biz.360es.cn:irp/connectors/agent-server.git --tags"
    if(result != 0){
        currentBuild.result = 'ABORTED'
        error('git tag error')
    }
}
