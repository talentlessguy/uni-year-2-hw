using System;
using System.Drawing;
using System.Windows.Forms;
using System.Collections.Generic;

namespace ThreeColumnLayout
{
    public partial class MainForm : Form
    {
        private List<(string ImagePath, string Caption)> imageList = new List<(string, string)>
        {
            ("daisy.jpg", "This is a daisy"),
            ("tulip.jpg", "This is a tulip"),
            ("heranium.jpg", "This is a heranium")
        };

        public MainForm()
        {
            InitializeUI();
            this.Size = new Size(800, 400);
        }


        private void InitializeUI()
        {
            // Create a TableLayoutPanel with 1 row and 3 columns
            TableLayoutPanel tableLayoutPanel = new TableLayoutPanel();
            tableLayoutPanel.Dock = DockStyle.Fill;
            tableLayoutPanel.RowCount = 1;
            tableLayoutPanel.ColumnCount = 3;

            // Create and add PictureBoxes with square images and buttons
            for (int i = 0; i <= 2; i++)
            {
                string imagePath = imageList[i].ImagePath;
                string caption = imageList[i].Caption;

                Image image = ResizeImageToSquare(Image.FromFile(imagePath), 256);


                // Create a FlowLayoutPanel for each image and button pair
                FlowLayoutPanel flowLayoutPanel = new FlowLayoutPanel();
                flowLayoutPanel.FlowDirection = FlowDirection.TopDown;
                flowLayoutPanel.AutoSize = true;

                // Create and add a PictureBox with a square image
                PictureBox pictureBox = new PictureBox();
                pictureBox.Size = new Size(256, 256);
                pictureBox.SizeMode = PictureBoxSizeMode.Zoom; // Maintain aspect ratio
                pictureBox.Image = image;
                flowLayoutPanel.Controls.Add(pictureBox);

                // Create and add a Button
                Button button = new Button();
                button.Text = $"Flower name {i}";
                button.Click += (sender, e) =>
                {
                    MessageBox.Show(caption, "Image Caption");
                };
                flowLayoutPanel.Controls.Add(button);

                // Add the FlowLayoutPanel to the TableLayoutPanel
                tableLayoutPanel.Controls.Add(flowLayoutPanel);
            }

            // Set column widths to evenly distribute space
            for (int i = 0; i < 3; i++)
            {
                tableLayoutPanel.ColumnStyles.Add(new ColumnStyle(SizeType.Percent, 33.33F));
            }

            // Add the TableLayoutPanel to the form
            Controls.Add(tableLayoutPanel);
        }

        private Image ResizeImageToSquare(Image img, int size)
        {
            // Create a new square bitmap
            Bitmap squareImage = new Bitmap(size, size);
            using (Graphics graphics = Graphics.FromImage(squareImage))
            {
                graphics.InterpolationMode = System.Drawing.Drawing2D.InterpolationMode.HighQualityBicubic;
                graphics.DrawImage(img, 0, 0, size, size);
            }
            return squareImage;
        }

        [STAThread]
        public static void Main(string[] args)
        {
            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);

            Application.Run(new MainForm());
        }
    }
}
