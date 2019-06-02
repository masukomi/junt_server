require 'airborne'

describe 'Junt' do
  before(:all) do
    ENV['JUNT_DB'] = "./test.db"
    # if ENV['JUNT_DB'].nil? 
    #   $stderr.puts("JUNT_DB env var must be specified.")
    #   exit 1
    # end
    # guarantee we're working with a fresh db
    `rm #{ENV['JUNT_DB']}`
    # start the junt server 
    @junt_server_pid = fork do
      exec "./junt"
    end
    puts "about to sleep..."
    # Process.detach(junt_server)
    # give the server a couple seconds to create the db
    # and start listening
    sleep(2) 


    @base_url = "http://localhost:8123/v1"
    @default_headers = {} # uses application/json content type by default
  end
  after(:all) do
    # Kill the server
    Process.kill("HUP", @junt_server_pid)
    # remove the test db
    `rm #{ENV['JUNT_DB']}`
  end
  it "should initially have no companies" do
    get "#{@base_url}/companies" # should return empty array
    expect(json_body.size).to(eq(0))
  end
  it "should be able to create a company" do
    new_company_json = {
      name: "test co",
      location: "somewhere nice",
      url: "https://example.com/test_co",
      note: "test co has 100% coverage"
    }
    post "#{@base_url}/companies", new_company_json, @default_headers
    expect_json_types(status: :string, id: :int)
    expect(json_body[:status]).to(eq("SUCCESS"))
  end
end
