using System;
using System.Numerics;
using System.Windows.Forms;

namespace FactorialCalculatorApp
{
    public class MainForm : Form
    {
        private Label resultLabel;
        private TextBox inputTextBox;
        private Button calculateButton;

        public MainForm()
        {
            InitializeComponents();
        }

        private void InitializeComponents()
        {
            resultLabel = new Label();
            inputTextBox = new TextBox();
            calculateButton = new Button();

            // Set properties for controls
            resultLabel.Location = new System.Drawing.Point(10, 80);
            resultLabel.Size = new System.Drawing.Size(300, 30);

            inputTextBox.Location = new System.Drawing.Point(10, 10);
            inputTextBox.Size = new System.Drawing.Size(100, 20);

            calculateButton.Location = new System.Drawing.Point(120, 10);
            calculateButton.Size = new System.Drawing.Size(80, 20);
            calculateButton.Text = "Calculate";
            calculateButton.Click += CalculateButton_Click;

            // Add controls to form
            Controls.Add(resultLabel);
            Controls.Add(inputTextBox);
            Controls.Add(calculateButton);

            // Set form properties
            Size = new System.Drawing.Size(300, 150);
            Text = "Factorial Calculator";
        }

        private void CalculateButton_Click(object sender, EventArgs e)
        {
            // Get user input
            if (BigInteger.TryParse(inputTextBox.Text, out BigInteger number) && number >= 0)
            {
                // Calculate factorial
                BigInteger factorial = CalculateFactorial(number);

                // Display result
                resultLabel.Text = $"Factorial of {number} is: {factorial}";
            }
            else
            {
                resultLabel.Text = "Invalid input. Please enter a non-negative integer.";
            }
        }

        private BigInteger CalculateFactorial(BigInteger n)
        {
            // Check if n is 0 or 1, return 1 in such cases
            if (n == 0 || n == 1)
            {
                return 1;
            }

            // Calculate factorial using a loop
            BigInteger result = 1;
            for (BigInteger i = 2; i <= n; i++)
            {
                result *= i;
            }

            return result;
        }

        [STAThread]
        public static void Main()
        {
            Application.Run(new MainForm());
        }
    }
}
