require 'chunky_png'

img = ChunkyPNG::Image.from_io(STDIN)

pixel_height = pixel_width = 20

cols = img.width / pixel_width
rows = img.height / pixel_height

class ColorArray
  def initialize
    @rgba = [0, 0, 0, 0]
  end

  def add(arr)
    @rgba[0] += arr[0]
    @rgba[1] += arr[1]
    @rgba[2] += arr[2]
    @rgba[3] += arr[3]
  end

  def avg(count)
    @rgba.map {|i| i / count }
  end
end

# @return [Integer] A measure of closeness, smaller means closer.
def closeness(a, b)
  parts(a).zip(parts(b)).map {|i,j| (i - j).abs }.inject(:+)
end

# @return [Array<Integer>] Parts of the colour
def parts(c)
  h = ChunkyPNG::Color
  [h.r(c), h.g(c), h.b(c), h.a(c)]
end

# @return [Integer] Average of the two colours +a+, +b+.
def average(a, b)
  r = ColorArray.new
  r.add(parts(a))
  r.add(parts(b))
  r.avg(2)
end

0.upto cols-1 do |col|
  0.upto rows-1 do |row|

    to = ColorArray.new
    bo = ColorArray.new
    le = ColorArray.new
    ri = ColorArray.new

    tc = bc = lc = rc = 0

    0.upto pixel_height-1 do |y|
      0.upto pixel_width-1 do |x|
        # Work out triangles seperately now
        real_y = pixel_height * row + y
        real_x = pixel_width * col + x

        # use center of block as origin, makes life easy
        y_origin = y - pixel_height / 2
        x_origin = x - pixel_width / 2

        # work out values to make life easier
        pixel = img[real_x, real_y]
        r, g, b, a = parts(pixel)

        # Top:
        if y_origin > x_origin && y_origin > -x_origin
          tc += 1
          to.add [r, g, b, a]

        # Right:
        elsif y_origin < x_origin && y_origin > -x_origin
          rc += 1
          ri.add [r, g, b, a]

        # Bottom:
        elsif y_origin < x_origin && y_origin < -x_origin
          bc += 1
          bo.add [r, g, b, a]

        # Left:
        elsif y_origin > x_origin && y_origin < -x_origin
          lc += 1
          le.add [r, g, b, a]
        end
      end
    end

    ato = ChunkyPNG::Color.rgba *to.avg(tc)
    ari = ChunkyPNG::Color.rgba *ri.avg(rc)
    abo = ChunkyPNG::Color.rgba *bo.avg(bc)
    ale = ChunkyPNG::Color.rgba *le.avg(lc)

    if closeness(ato, ari) > closeness(ato, ale)

      top_right = ChunkyPNG::Color.rgba *average(ato, ari)
      bottom_left = ChunkyPNG::Color.rgba *average(abo, ale)

      0.upto pixel_height-1 do |y|
        0.upto pixel_width-1 do |x|
          real_y = pixel_height * row + y
          real_x = pixel_width * col + x

          y_origin = y - pixel_height / 2
          x_origin = x - pixel_width / 2

          if y_origin > x_origin
            img[real_x, real_y] = top_right
          else
            img[real_x, real_y] = bottom_left
          end
        end
      end

    else

      top_left = ChunkyPNG::Color.rgba *average(ato, ale)
      bottom_right = ChunkyPNG::Color.rgba *average(abo, ari)

      0.upto pixel_height-1 do |y|
        0.upto pixel_width-1 do |x|
          real_y = pixel_height * row + y
          real_x = pixel_width * col + x

          y_origin = y - pixel_height / 2
          x_origin = x - pixel_width / 2

          if y_origin >= -x_origin
            img[real_x, real_y] = top_left
          else
            img[real_x, real_y] = bottom_right
          end
        end
      end
    end

  end
end

img.write(STDOUT)
