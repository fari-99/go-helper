package gohelper

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "strconv"
    "time"

    "github.com/jlaffaye/ftp"
    _ "github.com/joho/godotenv/autoload"
    "github.com/pkg/sftp"
    "golang.org/x/crypto/ssh"
)

const DefaultFTPPort = 21   // default ftp port
const DefaultFTPSPort = 990 // default ftps port
const DefaultSFTPPort = 22  // default sftp

type HelperFtp struct {
    FtpCredential FtpCredential `json:"ftp_credential"`
    isSsh         bool

    ftpFileLocation string
    ftpFilename     string
    isTest          bool
    isLocalTest     bool
}

type FtpCredential struct {
    FtpHost     string `json:"ftp_host"`
    FtpPort     string `json:"ftp_port"`
    FtpUser     string `json:"ftp_user"`
    FtpPassword string `json:"ftp_password"`
    SshUser     string `json:"ssh_user"`
    SshPassword string `json:"ssh_password"`
    SshKeyFile  string `json:"ssh_key_file"`
}

func BaseHelperFtp(isTest bool) *HelperFtp {
    localTest, _ := strconv.ParseBool(os.Getenv("FTP_LOCAL_TEST"))

    baseFtp := HelperFtp{
        isSsh:       false,
        isTest:      isTest,
        isLocalTest: localTest,
    }

    return &baseFtp
}

func (helper *HelperFtp) SetCredential(ftpCredential FtpCredential) *HelperFtp {
    if helper.isTest {
        ftpCredential = FtpCredential{
            FtpHost:     os.Getenv("FTP_TEST_HOST"),
            FtpPort:     os.Getenv("FTP_PORT"), // sftp port default 22
            SshUser:     os.Getenv("FTP_TEST_USERNAME"),
            SshPassword: os.Getenv("FTP_TEST_PASSWORD"),
            SshKeyFile:  os.Getenv("FTP_AUTH_FILE_LOCATION") + os.Getenv("FTP_TEST_AUTH_FILE"),
        }
    }

    if ftpCredential.SshUser != "" && ftpCredential.SshPassword != "" {
        helper.isSsh = true
    }

    helper.FtpCredential = ftpCredential
    return helper
}

func (helper *HelperFtp) SetFtpFile(targetFileLocation string, filename string) *HelperFtp {
    if helper.isTest {
        targetFileLocation = os.Getenv("FTP_TEST_LOCATION")
        filename = fmt.Sprintf("testing-%s", filename)
    }

    helper.ftpFileLocation = targetFileLocation
    helper.ftpFilename = filename
    return helper
}

func (helper *HelperFtp) SendFile(file interface{}) error {
    if helper.isTest {
        return helper.sftp(file)
    } else {
        ftpPort, err := strconv.ParseInt(helper.FtpCredential.FtpPort, 10, 64)
        if err != nil {
            return err
        }

        switch ftpPort {
        case DefaultFTPSPort, DefaultFTPPort:
            return helper.ftp(file)
        case DefaultSFTPPort:
            return helper.sftp(file)
        default:
            if ftpPort <= 0 {
                return fmt.Errorf("ftp port is invalid")
            }

            // using costum ports
            if !helper.isSsh {
                return helper.ftp(file)
            } else {
                return helper.sftp(file)
            }
        }
    }
}

func (helper *HelperFtp) getBuffer(file interface{}) (*bytes.Buffer, error) {
    newFile := bytes.NewBuffer(nil)
    switch fileType := file.(type) {
    case *os.File:
        _, _ = io.Copy(newFile, fileType)
    case string:
        openFile, err := os.Open(fileType)
        if err != nil {
            return nil, fmt.Errorf("error open file to get buffer, err := %s", err.Error())
        }

        defer openFile.Close()

        _, _ = io.Copy(newFile, openFile)
    case *bytes.Buffer:
        newFile = fileType
    default:
        return nil, fmt.Errorf("file not encoded, file type %T", file)
    }

    return newFile, nil
}

func (helper *HelperFtp) sftp(file interface{}) error {
    ftpCredential := helper.FtpCredential

    var authMethod []ssh.AuthMethod
    if ftpCredential.SshPassword != "" {
        authMethod = []ssh.AuthMethod{
            ssh.Password(ftpCredential.SshPassword),
        }
    } else if ftpCredential.SshKeyFile != "" {
        pemBytes, err := ioutil.ReadFile(ftpCredential.SshKeyFile)
        if err != nil {
            return fmt.Errorf("error read pem files.go, err := %s", err.Error())
        }

        signer, err := ssh.ParsePrivateKey(pemBytes)
        if err != nil {
            return fmt.Errorf("error parsing private key from pem fileLocation, err := %s", err.Error())
        }

        authMethod = []ssh.AuthMethod{
            ssh.PublicKeys(signer),
        }
    }

    sshConfig := &ssh.ClientConfig{
        User:            ftpCredential.SshUser,
        Auth:            authMethod,
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }

    address := fmt.Sprintf("%s:%s", ftpCredential.FtpHost, ftpCredential.FtpPort)
    sshConnection, err := ssh.Dial("tcp", address, sshConfig)
    if err != nil {
        return fmt.Errorf("error dial ssh authentication, err := %s", err.Error())
    }

    defer sshConnection.Close()

    sshSession, err := sshConnection.NewSession()
    if err != nil {
        return fmt.Errorf("error create ssh session, err := %s", err.Error())
    }

    defer sshSession.Close()

    sshSession.Stderr = os.Stdin
    sshSession.Stdin = os.Stdin
    sshSession.Stdout = os.Stdout

    sftpConnection, err := sftp.NewClient(sshConnection)
    if err != nil {
        return fmt.Errorf("failed to create sftp connection config pem fileLocation, err := %s", err.Error())
    }

    defer sftpConnection.Close()

    targetLocation := fmt.Sprintf("%s/%s", helper.ftpFileLocation, helper.ftpFilename)
    fileDestination, err := sftpConnection.Create(targetLocation)
    if err != nil {
        return fmt.Errorf("failed create destination fileLocation, err := %s", err.Error())
    }

    newFile, err := helper.getBuffer(file)
    if err != nil {
        return err
    }

    _, err = io.Copy(fileDestination, newFile)
    if err != nil {
        return fmt.Errorf("failed to copy data to fileLocation destination, err := %s", err.Error())
    }

    return nil
}

func (helper HelperFtp) ftp(file interface{}) error {
    ftpCredential := helper.FtpCredential

    ftpUser := ftpCredential.FtpUser
    ftpPassword := ftpCredential.FtpPassword

    address := fmt.Sprintf("%s:%s", ftpCredential.FtpHost, ftpCredential.FtpPort)
    clientFtp, err := ftp.Dial(address, ftp.DialWithTimeout(30*time.Second))
    if err != nil {
        return fmt.Errorf("error construct dial ftp host, err := %s", err.Error())
    }

    err = clientFtp.Login(ftpUser, ftpPassword)
    if err != nil {
        return fmt.Errorf("error login to ftp host, err := %s", err.Error())
    }

    targetLocation := fmt.Sprintf("%s/%s", helper.ftpFileLocation, helper.ftpFilename)

    newFile, err := helper.getBuffer(file)
    if err != nil {
        return err
    }

    err = clientFtp.Stor(targetLocation, newFile)
    if err != nil {
        return fmt.Errorf("error send fileLocation to ftp host, err := %s", err.Error())
    }

    return nil
}
