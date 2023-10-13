using System;
using System.Drawing;
using System.Windows.Forms;

namespace BMICalculator
{
    public class MainForm : Form
    {
        private Label weightLabel;
        private Label heightLabel;
        private TextBox weightTextBox;
        private TextBox heightTextBox;
        private Button calculateButton;
        private Label resultLabel;

        public MainForm()
        {
            InitializeComponents();
        }

        private void InitializeComponents()
        {
            weightLabel = new Label
            {
                Text = "Weight (pounds):",
                Location = new Point(20, 20),
                AutoSize = true
            };

            heightLabel = new Label
            {
                Text = "Height (inches):",
                Location = new Point(20, 60),
                AutoSize = true
            };

            weightTextBox = new TextBox
            {
                Location = new Point(150, 20),
                Size = new Size(100, 20)
            };

            heightTextBox = new TextBox
            {
                Location = new Point(150, 60),
                Size = new Size(100, 20)
            };

            calculateButton = new Button
            {
                Text = "Calculate BMI",
                Location = new Point(20, 100),
                Size = new Size(120, 30)
            };

            resultLabel = new Label
            {
                Text = "BMI: ",
                Location = new Point(20, 150),
                AutoSize = true
            };

            calculateButton.Click += CalculateButton_Click;

            Controls.Add(weightLabel);
            Controls.Add(heightLabel);
            Controls.Add(weightTextBox);
            Controls.Add(heightTextBox);
            Controls.Add(calculateButton);
            Controls.Add(resultLabel);
        }

        private void CalculateButton_Click(object sender, EventArgs e)
        {
            if (double.TryParse(weightTextBox.Text, out double weight) &&
                double.TryParse(heightTextBox.Text, out double height))
            {
                double bmi = CalculateBMI(weight, height);
                resultLabel.Text = $"BMI: {bmi:F2}";

                if (bmi < 18.5)
                    resultLabel.Text += " (Underweight)";
                else if (bmi >= 18.5 && bmi <= 25)
                    resultLabel.Text += " (Optimal Weight)";
                else
                    resultLabel.Text += " (Overweight)";
            }
            else
            {
                resultLabel.Text = "Invalid input. Please enter valid numbers.";
            }
        }

        private double CalculateBMI(double weight, double height)
        {
            // BMI formula: BMI = weight (kg) / (height (m))^2
            double heightInMeters = height * 0.0254; // Convert inches to meters
            return weight / (heightInMeters * heightInMeters);
        }

        [STAThread]
        public static void Main()
        {
            Application.EnableVisualStyles();
            Application.Run(new MainForm());
        }
    }
}
