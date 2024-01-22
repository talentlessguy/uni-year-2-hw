Option Explicit

Dim fso, startFolder, targetFolder, dict, folderCount
Set fso = CreateObject("Scripting.FileSystemObject")
Set dict = CreateObject("Scripting.Dictionary")
folderCount = 0

startFolder = WScript.Arguments(0)
targetFolder = WScript.Arguments(1)

SortFiles startFolder, targetFolder

WScript.Echo folderCount & " folders were created."

Dim key
For Each key In dict.Keys
    WScript.Echo "--------"
    WScript.Echo dict(key).Count & " file(s) were moved to folder"
    WScript.Echo key
    WScript.Echo "Files:"
    
    Dim fileKey
    For Each fileKey In dict(key).Keys
        WScript.Echo fileKey
    Next
Next

Function FormatDate(date)
    FormatDate = Year(date) & "\" & Year(date) & "-" & Right("0" & Month(date), 2) & "-" & Right("0" & Day(date), 2)
End Function

Sub SortFiles(folderPath, targetPath)
    Dim folder, file, fileDate, targetDateFolder, subFolder
    Set folder = fso.GetFolder(folderPath)

    For Each file In folder.Files
        Dim ext
        ext = LCase(fso.GetExtensionName(file))
        If ext = "jpg" Or ext = "jpeg" Or ext = "png" Then
            fileDate = FormatDate(file.DateLastModified)
            targetDateFolder = targetPath & "\" & fileDate
            If Not fso.FolderExists(targetPath & "\" & Year(file.DateLastModified)) Then 
                fso.CreateFolder(targetPath & "\" & Year(file.DateLastModified))
                folderCount = folderCount + 1
            End If
            If Not fso.FolderExists(targetDateFolder) Then 
                fso.CreateFolder(targetDateFolder)
                folderCount = folderCount + 1
            End If
            
            If Not dict.Exists(targetDateFolder) Then Set dict(targetDateFolder) = CreateObject("Scripting.Dictionary")
            dict(targetDateFolder).Add file.Name, True
            
            fso.MoveFile file.Path, targetDateFolder & "\"
        End If
    Next

    If folder.SubFolders.Count > 0 Then
        For Each subFolder In folder.SubFolders
            SortFiles subFolder.Path, targetPath
        Next
    End If
End Sub