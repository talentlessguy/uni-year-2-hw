using System;
using System.Collections.Generic;
using System.IO;
using System.Windows.Forms;

public class SpreadsheetEditor : Form
{
    private List<List<TextBox>> cells = new List<List<TextBox>>();

    public SpreadsheetEditor()
    {
        InitializeUI();
    }

    private void InitializeUI()
    {
        Button btnAppendRow = new Button
        {
            Text = "Append Row",
            Location = new System.Drawing.Point(10, 10)
        };
        btnAppendRow.Click += (sender, e) => AppendRow();

        Button btnAppendColumn = new Button
        {
            Text = "Append Column",
            Location = new System.Drawing.Point(120, 10)
        };
        btnAppendColumn.Click += (sender, e) => AppendColumn();

        Button btnSaveCSV = new Button
        {
            Text = "Save to CSV",
            Location = new System.Drawing.Point(230, 10)
        };
        btnSaveCSV.Click += (sender, e) => SaveToCSV();

        Controls.Add(btnAppendRow);
        Controls.Add(btnAppendColumn);
        Controls.Add(btnSaveCSV);

        InitializeGrid(5, 5);
    }

    private void InitializeGrid(int rows, int columns)
    {
        for (int i = 0; i < rows; i++)
        {
            List<TextBox> row = new List<TextBox>();
            for (int j = 0; j < columns; j++)
            {
                TextBox cell = new TextBox
                {
                    Location = new System.Drawing.Point(j * 50 + 10, i * 20 + 50),
                    Width = 40
                };
                row.Add(cell);
                Controls.Add(cell);
            }
            cells.Add(row);
        }
    }

    private void AppendRow()
    {
        int columns = cells[0].Count;
        List<TextBox> newRow = new List<TextBox>();

        for (int j = 0; j < columns; j++)
        {
            TextBox cell = new TextBox
            {
                Location = new System.Drawing.Point(j * 50 + 10, cells.Count * 20 + 50),
                Width = 40
            };
            newRow.Add(cell);
            Controls.Add(cell);
        }
        cells.Add(newRow);
    }

    private void AppendColumn()
    {
        int rows = cells.Count;
        int newColumnIndex = cells[0].Count;

        for (int i = 0; i < rows; i++)
        {
            TextBox cell = new TextBox
            {
                Location = new System.Drawing.Point(newColumnIndex * 50 + 10, i * 20 + 50),
                Width = 40
            };
            cells[i].Add(cell);
            Controls.Add(cell);
        }
    }

    private void SaveToCSV()
    {
        using (StreamWriter writer = new StreamWriter("spreadsheet.csv"))
        {
            foreach (var row in cells)
            {
                string line = string.Join(",", row.ConvertAll(cell => cell.Text));
                writer.WriteLine(line);
            }
        }
        MessageBox.Show("CSV file saved successfully!");
    }

    [STAThread]
    public static void Main()
    {
        Application.EnableVisualStyles();
        Application.SetCompatibleTextRenderingDefault(false);
        Application.Run(new SpreadsheetEditor());
    }
}