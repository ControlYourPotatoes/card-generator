// Function to extract transcript text from the JSON data
function extractTranscript(jsonData) {
    // Parse the JSON if it's a string
    let data = typeof jsonData === 'string' ? JSON.parse(jsonData) : jsonData;
    
    // Check if data is an array or if it has a data property that is an array
    let transcriptEntries = Array.isArray(data) ? data : data.data;
    
    if (!transcriptEntries || !Array.isArray(transcriptEntries)) {
      return "Error: Could not find transcript entries in the provided JSON.";
    }
    
    // Extract text from each entry's lines
    let transcript = [];
    
    transcriptEntries.forEach(entry => {
      if (entry.lines && Array.isArray(entry.lines)) {
        entry.lines.forEach(line => {
          if (line.text) {
            transcript.push(line.text);
          }
        });
      }
    });
    
    // Join the transcript lines with a space to form complete sentences
    return transcript.join(" ");
  }
  
  // Node.js implementation to process your file
  const fs = require('fs');
  const path = require('path');
  
  // Get file path from command line arguments
  // If no argument is provided, use a default value
  const jsonFilePath = process.argv[2] || 'your-transcript-file.json';
  
  // Create output filename based on input filename
  const inputFileName = path.basename(jsonFilePath, path.extname(jsonFilePath));
  const outputFilePath = `${inputFileName}-transcript.txt`;
  
  try {
    // Read and parse the JSON file
    const jsonData = JSON.parse(fs.readFileSync(jsonFilePath, 'utf8'));
    
    // Extract the transcript
    const transcript = extractTranscript(jsonData);
    
    // Save the transcript to a text file
    fs.writeFileSync(outputFilePath, transcript, 'utf8');
    
    console.log('Transcript extraction complete!');
    console.log(`Saved to: ${outputFilePath}`);
    
    // Also display the first 150 characters of the transcript
    console.log('\nPreview:');
    console.log(transcript.substring(0, 150) + '...');
  } catch (error) {
    console.error('Error processing the transcript:', error.message);
  }