using System;
using System.Windows.Forms;
using System.Drawing;

public class TicketSalesProgram : Form
{
    private Label ticketsSoldLabel;
    private Label revenueGeneratedLabel;

    private Label classALabel;
    private Label classBLabel;
    private Label classCLabel;

    private TextBox classATicketsTextBox;
    private TextBox classBTicketsTextBox;
    private TextBox classCTicketsTextBox;

    private TextBox classARevenueTextBox;
    private TextBox classBRevenueTextBox;
    private TextBox classCRevenueTextBox;
    private TextBox totalRevenueTextBox;

    private Button calculateButton;
    private Button clearButton;
    private Button exitButton;

    public TicketSalesProgram()
    {
        this.Text = "Ticket Sales Program";
        this.Size = new System.Drawing.Size(800, 400);
        this.StartPosition = FormStartPosition.CenterScreen;

        InitializeUI();
    }

    private void InitializeUI()
    {
        ticketsSoldLabel = new Label();
        ticketsSoldLabel.Text = "Tickets Sold";
        ticketsSoldLabel.Location = new System.Drawing.Point(30, 20);
        this.Controls.Add(ticketsSoldLabel);

        revenueGeneratedLabel = new Label();
        revenueGeneratedLabel.Text = "Revenue Generated";
        revenueGeneratedLabel.Location = new System.Drawing.Point(280, 20);
        this.Controls.Add(revenueGeneratedLabel);

        classALabel = new Label();
        classALabel.Text = "Class A:";
        classALabel.Location = new System.Drawing.Point(30, 50);
        this.Controls.Add(classALabel);

        classBLabel = new Label();
        classBLabel.Text = "Class B:";
        classBLabel.Location = new System.Drawing.Point(30, 80);
        this.Controls.Add(classBLabel);

        classCLabel = new Label();
        classCLabel.Text = "Class C:";
        classCLabel.Location = new System.Drawing.Point(30, 110);
        this.Controls.Add(classCLabel);

        // Textboxes for ticket quantities
        classATicketsTextBox = new TextBox();
        classATicketsTextBox.Location = new System.Drawing.Point(120, 50);
        this.Controls.Add(classATicketsTextBox);

        classBTicketsTextBox = new TextBox();
        classBTicketsTextBox.Location = new System.Drawing.Point(120, 80);
        this.Controls.Add(classBTicketsTextBox);

        classCTicketsTextBox = new TextBox();
        classCTicketsTextBox.Location = new System.Drawing.Point(120, 110);
        this.Controls.Add(classCTicketsTextBox);

        classARevenueTextBox = new TextBox();
        classARevenueTextBox.Location = new System.Drawing.Point(300, 50);
        classARevenueTextBox.ReadOnly = true;
        this.Controls.Add(classARevenueTextBox);

        classBRevenueTextBox = new TextBox();
        classBRevenueTextBox.Location = new System.Drawing.Point(300, 80);
        classBRevenueTextBox.ReadOnly = true;
        this.Controls.Add(classBRevenueTextBox);

        classCRevenueTextBox = new TextBox();
        classCRevenueTextBox.Location = new System.Drawing.Point(300, 110);
        classCRevenueTextBox.ReadOnly = true;
        this.Controls.Add(classCRevenueTextBox);

        totalRevenueTextBox = new TextBox();
        totalRevenueTextBox.Location = new System.Drawing.Point(300, 140);
        totalRevenueTextBox.ReadOnly = true;
        this.Controls.Add(totalRevenueTextBox);

        calculateButton = new Button();
        calculateButton.Text = "Calculate Revenue";
        calculateButton.Location = new System.Drawing.Point(30, 170);
        calculateButton.Click += new EventHandler(CalculateButton_Click);
        this.Controls.Add(calculateButton);

        clearButton = new Button();
        clearButton.Text = "Clear";
        clearButton.Location = new System.Drawing.Point(140, 170);
        clearButton.Click += new EventHandler(ClearButton_Click);
        this.Controls.Add(clearButton);

        exitButton = new Button();
        exitButton.Text = "Exit";
        exitButton.Location = new System.Drawing.Point(250, 170);
        exitButton.Click += new EventHandler(ExitButton_Click);
        this.Controls.Add(exitButton);
    }

    private void CalculateButton_Click(object sender, EventArgs e)
    {
        int classAQuantity = string.IsNullOrEmpty(classATicketsTextBox.Text) ? 0 : int.Parse(classATicketsTextBox.Text);
        int classBQuantity = string.IsNullOrEmpty(classBTicketsTextBox.Text) ? 0 : int.Parse(classBTicketsTextBox.Text);
        int classCQuantity = string.IsNullOrEmpty(classCTicketsTextBox.Text) ? 0 : int.Parse(classCTicketsTextBox.Text);

        double classARevenue = classAQuantity * 15;
        double classBRevenue = classBQuantity * 12;
        double classCRevenue = classCQuantity * 9;

        classARevenueTextBox.Text = classARevenue.ToString("C0");
        classBRevenueTextBox.Text = classBRevenue.ToString("C0");
        classCRevenueTextBox.Text = classCRevenue.ToString("C0");

        double totalRevenue = classARevenue + classBRevenue + classCRevenue;
        totalRevenueTextBox.Text = totalRevenue.ToString("C0");
    }

    private void ClearButton_Click(object sender, EventArgs e)
    {
        classATicketsTextBox.Clear();
        classBTicketsTextBox.Clear();
        classCTicketsTextBox.Clear();
        classARevenueTextBox.Clear();
        classBRevenueTextBox.Clear();
        classCRevenueTextBox.Clear();
        totalRevenueTextBox.Clear();
    }

    private void ExitButton_Click(object sender, EventArgs e)
    {
        this.Close();
    }

    public static void Main()
    {
        
        Application.Run(new TicketSalesProgram());
    }
}
