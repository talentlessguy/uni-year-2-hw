Option Explicit

Sub CheckZipfsLaw(filePath, intUserSpecifiedNum)
    Dim objFSO, objTextFile
    Dim strText, arrWords, word, dict, dictApostrophe
    Dim i, intTotalWords, normalizedWord
    Dim arrKeys, arrItems
    Dim wordForms
    Dim regex

    Set objFSO = CreateObject("Scripting.FileSystemObject")
    Set objTextFile = objFSO.OpenTextFile(filePath, 1)

    strText = LCase(objTextFile.ReadAll)
    objTextFile.Close

    strText = Replace(Replace(strText, Chr(13), " "), Chr(10), " ")
    Set regex = New RegExp
    regex.Pattern = "[^a-zA-Z\s']"
    regex.Global = True
    strText = regex.Replace(strText, "")

    Set wordForms = CreateObject("Scripting.Dictionary")

    wordForms.Add "the", Array("this", "that", "these", "those")
    wordForms.Add "be", Array("is", "are", "am", "was", "were")
    wordForms.Add "to", Array("toward", "unto", "until", "till")
    wordForms.Add "of", Array("off", "off of")
    wordForms.Add "and", Array("also", "plus", "as well as")
    wordForms.Add "a", Array("an")
    wordForms.Add "in", Array("inside", "within", "into")
    wordForms.Add "have", Array("has", "had")
    wordForms.Add "it", Array("its")
    wordForms.Add "not", Array("don't", "doesn't", "didn't")
    wordForms.Add "he", Array("him", "his")
    wordForms.Add "that", Array("those", "this")
    wordForms.Add "you", Array("your", "yours")
    wordForms.Add "do", Array("does", "did")
    wordForms.Add "are", Array("is", "am", "was", "were")
    wordForms.Add "this", Array("that")
    wordForms.Add "but", Array("however", "nevertheless")
    wordForms.Add "on", Array("upon", "onto")
    wordForms.Add "with", Array("along with", "together with")
    wordForms.Add "i", Array("me", "my", "mine")
    wordForms.Add "at", Array("around", "about")
    wordForms.Add "by", Array("beside", "near", "next to")
    wordForms.Add "they", Array("them", "their", "theirs")
    wordForms.Add "we", Array("us", "our", "ours")
    wordForms.Add "say", Array("said", "says")
    wordForms.Add "she", Array("hers", "her")
    wordForms.Add "for", Array("four", "fore", "before")
    wordForms.Add "good", Array("great", "excellent", "positive", "beneficial")
    wordForms.Add "time", Array("moment", "occasion", "event")
    wordForms.Add "work", Array("job", "employment", "task", "labor")

    arrWords = Split(strText)

    For i = 0 To UBound(arrWords)
        For Each word In wordForms.Keys
            For Each normalizedWord In wordForms.Item(word)
                If arrWords(i) = normalizedWord Then
                    arrWords(i) = word
                End If
            Next
        Next
    Next

    Set dict = CreateObject("Scripting.Dictionary")
    Set dictApostrophe = CreateObject("Scripting.Dictionary")

    For Each word In arrWords
        If word <> "" Then
            If InStr(word, "'") > 0 Then
                If dictApostrophe.Exists(word) Then
                    dictApostrophe.Item(word) = dictApostrophe.Item(word) + 1
                Else
                    dictApostrophe.Add word, 1
                End If
            Else
                If dict.Exists(word) Then
                    dict.Item(word) = dict.Item(word) + 1
                Else
                    dict.Add word, 1
                End If
            End If
        End If
    Next

    intTotalWords = UBound(arrWords) + 1

    WScript.Echo "The most popular words in " & filePath & " are: " & vbCrLf
    arrKeys = dict.Keys
    arrItems = dict.Items
    Call BubbleSort(arrKeys, arrItems)
    For i = 0 To intUserSpecifiedNum - 1
        WScript.Echo arrKeys(i) & " " & arrItems(i) & " " & intTotalWords / (i + 1)
    Next

    WScript.Echo vbCrLf & "The most popular still remaining short forms in " & filePath & " are: " & vbCrLf
    arrKeys = dictApostrophe.Keys
    arrItems = dictApostrophe.Items
    Call BubbleSort(arrKeys, arrItems)
    For i = 0 To intUserSpecifiedNum - 1
        WScript.Echo arrKeys(i) & " " & arrItems(i) & " " & intTotalWords / (i + 1)
    Next
End Sub

Sub BubbleSort(arr1, arr2)
    Dim i, j, temp1, temp2
    For i = UBound(arr1) - 1 To 0 Step -1
        For j = 0 To i
            If arr2(j) < arr2(j + 1) Then
                temp1 = arr1(j)
                temp2 = arr2(j)
                arr1(j) = arr1(j + 1)
                arr2(j) = arr2(j + 1)
                arr1(j + 1) = temp1
                arr2(j + 1) = temp2
            End If
        Next
    Next
End Sub

Dim filePath, intUserSpecifiedNum

If WScript.Arguments.Count >= 2 Then
    filePath = WScript.Arguments(0)
    intUserSpecifiedNum = WScript.Arguments(1)
    Call CheckZipfsLaw(filePath, intUserSpecifiedNum)
Else
    WScript.Echo "Please provide a file path and a number as arguments."
End If
